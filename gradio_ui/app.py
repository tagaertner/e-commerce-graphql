from gradio_ui.interface import build_interface

"""
app.py
------
Main Gradio application entry point.
Loads components and assembles the UI.
"""






if __name__ == "__main__":
    app = build_interface()
    app.launch(server_name="0.0.0.0", server_port=4103)
    # TODO need to change ports