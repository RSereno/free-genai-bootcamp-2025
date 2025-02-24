import os

# Add these lines at the top of your file, before any other imports
os.environ["OTEL_SDK_DISABLED"] = "true"
os.environ["OTEL_TRACES_EXPORTER"] = "none"

from comps import ServiceOrchestrator, MicroService, ServiceType, ServiceRoleType
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatCompletionResponseChoice,
    ChatMessage,
    UsageInfo
)
from fastapi import Request
from fastapi.responses import StreamingResponse
from comps.cores.mega.utils import handle_message
from comps.cores.proto.docarray import LLMParams, RerankerParms, RetrieverParms



EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 6000)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "0.0.0.0")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 11434)


class ExampleService:
    def __init__(self, host="0.0.0.0", port=8000):
        self.host = host
        self.port = port
        self.endpoint = "/v1/example-service"
        self.megaservice = ServiceOrchestrator()

    def add_remote_service(self):
        # embedding = MicroService(
        #     name="embedding",
        #     host=EMBEDDING_SERVICE_HOST_IP,
        #     port=EMBEDDING_SERVICE_PORT,
        #     endpoint="/v1/embeddings",
        #     use_remote_service=True,
        #     service_type=ServiceType.EMBEDDING,
        # )
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/v1/chat/completions",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        # self.megaservice.add(embedding).add(llm)
        # self.megaservice.flow_to(embedding, llm)
        self.megaservice.add(llm)
    
    async def handle_request(self, request: Request):
        data = await request.json()
        stream_opt = data.get("stream", True)
        chat_request = ChatCompletionRequest.model_validate(data)
        # prompt = handle_message(chat_request.messages)
        parameters = LLMParams(
            max_tokens=chat_request.max_tokens if chat_request.max_tokens else 1024,
            top_k=chat_request.top_k if chat_request.top_k else 10,
            top_p=chat_request.top_p if chat_request.top_p else 0.95,
            temperature=chat_request.temperature if chat_request.temperature else 0.01,
            frequency_penalty=chat_request.frequency_penalty if chat_request.frequency_penalty else 0.0,
            presence_penalty=chat_request.presence_penalty if chat_request.presence_penalty else 0.0,
            repetition_penalty=chat_request.repetition_penalty if chat_request.repetition_penalty else 1.03,
            stream=stream_opt,
            chat_template=chat_request.chat_template if chat_request.chat_template else None,
        )
        # retriever_parameters = RetrieverParms(
        #     search_type=chat_request.search_type if chat_request.search_type else "similarity",
        #     k=chat_request.k if chat_request.k else 4,
        #     distance_threshold=chat_request.distance_threshold if chat_request.distance_threshold else None,
        #     fetch_k=chat_request.fetch_k if chat_request.fetch_k else 20,
        #     lambda_mult=chat_request.lambda_mult if chat_request.lambda_mult else 0.5,
        #     score_threshold=chat_request.score_threshold if chat_request.score_threshold else 0.2,
        # )
        # reranker_parameters = RerankerParms(
        #     top_n=chat_request.top_n if chat_request.top_n else 1,
        # )
        initial_inputs={
            "messages": chat_request.messages,
        }
        result_dict, runtime_graph = await self.megaservice.schedule(
            #initial_inputs={"text": prompt},
            initial_inputs=initial_inputs,
            llm_parameters=parameters
            # ,
            # retriever_parameters=retriever_parameters,
            # reranker_parameters=reranker_parameters,
        )
        for node, response in result_dict.items():
            if isinstance(response, StreamingResponse):
                return response
        last_node = runtime_graph.all_leaves()[-1]

        print("last_node:\n",last_node)
        response = result_dict[last_node]["text"]
        choices = []
        usage = UsageInfo()
        choices.append(
            ChatCompletionResponseChoice(
                index=0,
                message=ChatMessage(role="assistant", content=response),
                finish_reason="stop",
            )
        )
        return ChatCompletionResponse(model="chatqna", choices=choices, usage=usage)

    def start(self):

        self.service = MicroService(
            self.__class__.__name__,
            service_role=ServiceRoleType.MEGASERVICE,
            host=self.host,
            port=self.port,
            endpoint=self.endpoint,
            input_datatype=ChatCompletionRequest,
            output_datatype=ChatCompletionResponse,
        )

        self.service.add_route(self.endpoint, self.handle_request, methods=["POST"])

        self.service.start()
    
example = ExampleService()
example.add_remote_service()
example.start()