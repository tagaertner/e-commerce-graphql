"""
interface.py
Outline for Gradio tabs + UI layout.
"""


import gradio as gr  # type: ignore
from handlers import (
    handle_create_user, handle_list_users, handle_get_user,
    handle_update_user, handle_delete_user,

    handle_add_product, handle_list_products, handle_get_product,
    handle_update_product, handle_delete_product,

    handle_create_order, handle_get_orders_for_user,
    handle_get_order, handle_update_order, handle_delete_order
)

def build_interface():
    with gr.Blocks() as demo:

        gr.Markdown("# 🛍️ E-Commerce Portal Showcase")

        # =======================
        # USERS TAB
        # =======================
        with gr.Tab("Users"):

            gr.Markdown("## 👤 Create Account")

            name = gr.Textbox(label="Name")
            email = gr.Textbox(label="Email")
            password = gr.Textbox(label="Password", type="password")
            active = gr.Checkbox(label="Active", value=True)

            create_btn = gr.Button("Create Account")
            create_output = gr.Textbox(label="Result", interactive=False)

            # (Wire handler later)
            create_btn.click(
                fn=handle_create_user,
                inputs=[name, email, password, active],
                outputs=[create_output, name, email, password, active]
            )

            gr.Markdown("---")
            gr.Markdown("## 🔐 Admin Actions")

            # Staying empty for now — placeholder only
            gr.Markdown("*Admin features coming later.*")

        # =======================
        # PRODUCTS TAB
        # =======================
        with gr.Tab("Products"):

            gr.Markdown("## 🛍️ Customer View")
            
            # Pagination state
            cursor_state = gr.State(value=None)
            
            with gr.Row():
                page_size = gr.Number(
                    value=10,
                    minimum=1,
                    maximum=50,
                    step=1,
                    label="Products per page",
                )
                load_btn = gr.Button("Load Products 🔄")
                next_btn = gr.Button("Next Page 📃") 
                
            products_table = gr.Dataframe(
                headers=["ID", "Name", "Price"],
                interactive=False,
                wrap=True,
                label="Product List",
            ) 
            
            gr.Markdown("### 🔍 Product Details")
            
            product_id_input = gr.Textbox(
                label="Product ID",
                placeholder="Paste or type a product ID from the table above",
            )
            product_details_btn =gr.Button("View Product Details")
            product_details_output = gr.Textbox(
                label="Product Details",
                lines=8,
                interactive=False,
            )
            
            #  === EVENT HANDLERS === 
            load_btn.click(
                fn=lambda first: handle_list_products(None, first),
                inputs=[page_size],
                outputs=[products_table, cursor_state]
            )
            
            next_btn.click(
                fn=handle_list_products,
                inputs=[cursor_state, page_size],
                 outputs=[products_table, cursor_state],
            )
            
            product_details_btn.click(
                fn=handle_get_product,
                inputs=[product_id_input],
                outputs=[product_details_output],
            )
            gr.Markdown("*Product listing UI will go here*")

            gr.Markdown("---")
            gr.Markdown("## 🔐 Admin Management")
            gr.Markdown("*Admin product tools coming later.*")

        # =======================
        # ORDERS TAB
        # =======================
        with gr.Tab("Orders"):

            gr.Markdown("## 🛒 Customer Orders")
            gr.Markdown("*Order creation UI will go here*")

            gr.Markdown("---")
            gr.Markdown("## 🔐 Admin Order Tools")
            gr.Markdown("*Admin order tools coming later.*")

    return demo