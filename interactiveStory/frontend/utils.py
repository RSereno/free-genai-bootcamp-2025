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
                st.session_state.scene_id = option["nextScene"]
                st.session_state.points += option.get("points", 0)
                st.rerun()
        
        # Display hint if available in the scene or options
        hint = scene.get("hint")
        if not hint:  # Check if hint is inside options
            for option in options:
                if "hint" in option:
                    hint = option["hint"]
                    break
        if hint:
            st.info(f"Hint: {hint}")
        
        # Display current points
        st.sidebar.write(f"🌟 Points: {st.session_state.points}")
    else:
        st.error("Erro: Nenhuma cena encontrada. Verifique o backend. 🚨")
