from interface import build_interface

"""
app.py
------
Main Gradio application entry point.
Loads components and assembles the UI.
"""






if __name__ == "__main__":
    demo = build_interface()
    demo.launch(server_name="0.0.0.0", server_port=4004)
    # TODO need to change ports