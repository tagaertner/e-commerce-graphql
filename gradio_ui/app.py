from interface import build_interface
import os

"""
app.py
------
Main Gradio application entry point.
Loads components and assembles the UI.
"""


    
if __name__ == "__main__":
    port = int(os.environ.get("PORT", 10000))
    demo = build_interface()
    
    print(f"🚀 Starting Gradio on port {port}")
    
    demo.launch(
        server_name="0.0.0.0", 
        server_port=port,
        root_path="/",
        share=False
        )
 