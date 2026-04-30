"""
interface.py
Outline for Gradio tabs + UI layout.
"""

import gradio as gr  # type: ignore
from handlers import (
    handle_create_user,
    handle_list_products,
    handle_get_product,
    handle_get_orders_for_user,
    handle_add_to_basket,
    handle_create_order_from_basket,
)


def select_product_row(products, event: gr.SelectData):
    if not products:
        return ""

    row, _ = event.index

    if row < 0 or row >= len(products):
        return ""

    return products[row][0]


def build_interface():
    with gr.Blocks() as demo:
        gr.Markdown("# 🛍️ E-Commerce Portal Showcase")

        cursor_state = gr.State(value=None)
        products_state = gr.State(value=[])
        selected_product_id = gr.State(value="")
        basket_state = gr.State(value=[])

        with gr.Tab("Users"):
            gr.Markdown("## 👤 Create Account")

            name = gr.Textbox(label="Name")
            email = gr.Textbox(label="Email")
            password = gr.Textbox(label="Password", type="password")
            active = gr.Checkbox(label="Active", value=True)

            create_user_btn = gr.Button("Create Account")
            create_user_output = gr.Textbox(label="Result", interactive=False)

        with gr.Tab("Products"):
            gr.Markdown("## 🛍️ Customer View")

            with gr.Row():
                page_size = gr.Number(value=10, minimum=1, maximum=50, step=1, label="Products per page")
                load_products_btn = gr.Button("Load Products 🔄")
                next_products_btn = gr.Button("Next Page ⏭️")

            products_table = gr.Dataframe(
                headers=["ID", "Name", "Price"],
                interactive=False,
                wrap=True,
                label="Product List",
            )

            gr.Markdown("### 🔍 Product Details")

            product_details_output = gr.Textbox(label="Product Details", lines=8, interactive=False)

            gr.Markdown("### 🧺 Add Selected Product to Basket")

            basket_quantity = gr.Number(label="Quantity", value=1, minimum=1, step=1)
            add_to_basket_btn = gr.Button("Add to Basket")
            basket_message = gr.Textbox(label="Basket Message", interactive=False)

        with gr.Tab("Orders"):
            gr.Markdown("## 🧺 Basket")

            basket_table = gr.Dataframe(
                headers=["Product ID", "Name", "Price", "Quantity", "Line Total"],
                label="Basket",
                interactive=False,
            )

            gr.Markdown("---")
            gr.Markdown("## 📦 View Orders for User")

            order_user_id = gr.Textbox(label="User ID", interactive=False)
            view_orders_btn = gr.Button("View Orders")

            orders_output = gr.Textbox(
                label="Orders",
                lines=12,
                interactive=False,
            )

    

        # =======================
        # EVENT HANDLERS
        # =======================

        create_user_btn.click(
            fn=handle_create_user,
            inputs=[name, email, password, active],
            outputs=[create_user_output, name, email, password, active],
        )

        load_products_btn.click(
            fn=lambda first: handle_list_products(None, first),
            inputs=[page_size],
            outputs=[products_table, cursor_state, products_state],
        )

        next_products_btn.click(
            fn=handle_list_products,
            inputs=[cursor_state, page_size],
            outputs=[products_table, cursor_state, products_state],
        )

        products_table.select(
            fn=select_product_row,
            inputs=[products_state],
            outputs=[selected_product_id],
        )

        selected_product_id.change(
            fn=handle_get_product,
            inputs=[selected_product_id],
            outputs=[product_details_output],
        )

        add_to_basket_btn.click(
            fn=handle_add_to_basket,
            inputs=[selected_product_id, basket_quantity, basket_state],
            outputs=[basket_state, basket_table, basket_message],
        )

        view_orders_btn.click(
            fn=handle_get_orders_for_user,
            inputs=[order_user_id],
            outputs=[orders_output],
        )

    return demo