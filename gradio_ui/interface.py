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
    handle_create_order, handle_list_orders, handle_list_orders_for_user,
    handle_get_order, handle_update_order, handle_delete_order
)

def build_interface():
    with gr.Blocks() as demo:

        gr.Markdown("# 🛍️ E-Commerce Portal Showcase")

        # =======================
        # USERS TAB
        # =======================
        with gr.Tab("Users"):

            gr.Markdown("## 👤 Customer Actions")

            # --- Customer Features ---
            # register
            # login
            # view my account
            # update my account

            gr.Markdown("---")

            gr.Markdown("## 🔐 Admin Actions")

            # --- Admin Features ---
            # list all users
            # get user by id
            # delete user


        # =======================
        # PRODUCTS TAB
        # =======================
        with gr.Tab("Products"):

            gr.Markdown("## 🛍️ Customer View")
            # list products
            # view product details

            gr.Markdown("---")

            gr.Markdown("## 🔐 Admin Management")
            # add product
            # update product
            # delete product


        # =======================
        # ORDERS TAB
        # =======================
        with gr.Tab("Orders"):

            gr.Markdown("## 🛒 Customer Orders")
            # create order
            # view my orders

            gr.Markdown("---")

            gr.Markdown("## 🔐 Admin Order Tools")
            # list all orders
            # get orders by user
            # update any order
            # delete order

    return demo