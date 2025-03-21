import streamlit as st

def display_scene(scene):
    """Display the current scene and handle user interactions."""
    if "error" in scene:
        st.error("Nenhuma cena encontrada.")
    elif "prompt" in scene:
        st.write(scene["prompt"])
        options = scene.get("options", [])
        
        for option in options:
            # Add a unique key to each button
            if st.button(option["text"], key=f"{option['text']}_{option['nextScene']}"):
                if "hint" in option:  # Check if the option has a hint
                    st.error("Wrong choice! Try again.")
                    st.info(f"Hint: {option['hint']}")
                    st.stop()  # Halt execution to ensure the error is displayed
                st.session_state.scene_id = option["nextScene"]
                st.session_state.points += option.get("points", 0)
                st.rerun()
        
        # Display current points
        st.sidebar.write(f"🌟 Points: {st.session_state.points}")
    else:
        st.error("Erro: Nenhuma cena encontrada. Verifique o backend. 🚨")
