"""
handlers.py
Transform UI inputs → GraphQL client calls.
"""
import gradio as gr
from datetime import datetime
from graphql_client import (
    # User client functions
    create_user, get_users, get_user_by_id, update_user, delete_user,

    # Product client functions
    get_products, get_product_by_id, create_product, update_product, delete_product,
    restock_product, set_product_availability,

    # Order client functions
    create_order, get_orders, get_orders_for_user, get_order_by_id,
    update_order, delete_order, set_order_status, change_order_quantity
)

"""
===========================
 USER HANDLERS (Admin + Customer)
===========================
"""

def create_user(name, email, password, role, active):
    variables = {
        "input": {
            "name": name,
            "email": email,
            "password": password,
            "role": role,
            "active": active,
        }
    }
    
    result = create_user(variables["input"])
    
    if result is None:
        return "❌ Error: No response from GraphQL server", "", "", "", False

    if "errors" in result:
        return f"❌ Error: {result['errors'][0]['message']}", "", "", "", False

    if "createUser" not in result["data"]:
        return f"❌ Error: Unexpected response format: {result}", "", "", "", False

    user = result["data"]["createUser"]
    message = (
        f"✅ User Created!\n\n"
        f"ID: {user['id']}\n"
        f"Name: {user['name']}\n"
        f"Email: {user['email']}\n"
        f"Role: {user['role']}\n"
        f"Active: {user['active']}"
    )

    # Return new message + clear form fields for Gradio
    return message, "", "", "", False

# ADMIN
def list_users(event: gr.SelectData):
    # event.value → the selected cell value
    # event.index → (row_index, col_index)

    _, col = event.index

    # Only allow clicking on the first column (user ID)
    if col != 0:
        return "", "❌ You can only select by clicking the User ID column."

    user_id = event.value
    return user_id, f"✅ Selected User ID: {user_id}"

def get_user(user_id):
    pass  # ADMIN

def update_user(user_id, name, email):
    pass

def delete_user(user_id):
    pass



"""
===========================
 PRODUCT HANDLERS (Admin + Customer)
===========================
"""

def add_product(name, price, description, inventory):
    pass  # ADMIN

def list_products():
    pass  # CUSTOMER

def get_product(product_id):
    pass  # CUSTOMER

def update_product(product_id, name, price, description, inventory):
    pass  # ADMIN

def delete_product(product_id):
    pass  # ADMIN

def restock_product(product_id, quantity):
    pass  # ADMIN

def set_product_availability(product_id, available):
    pass  # ADMIN



"""
===========================
 ORDER HANDLERS (Admin + Customer)
===========================
"""

def create_order(user_id, product_ids, quantity, total_price):
    pass  # CUSTOMER

def list_orders():
    pass  # ADMIN

def list_orders_for_user(user_id):
    pass  # CUSTOMER

def get_order(order_id):
    pass  # ADMIN

def update_order(order_id, quantity=None, total_price=None, status=None):
    pass  # ADMIN

def delete_order(order_id, user_id):
    pass  # ADMIN

def set_order_status(order_id, status):
    pass  # ADMIN

def change_order_quantity(order_id, quantity):
    pass  # ADMIN