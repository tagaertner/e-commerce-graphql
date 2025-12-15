"""
handlers.py
Transform UI inputs → GraphQL client calls.
"""
import gradio as gr
from datetime import datetime
from graphql_client import (
    # Users
    create_user as gql_create_user,
    get_users as gql_get_users,
    get_user_by_id as gql_get_user_by_id,
    update_user as gql_update_user,
    delete_user as gql_delete_user,

    # Products

    get_products_cursor as gql_get_products_cursor,
    get_product_by_id as gql_get_product_by_id,
    create_product as gql_create_product,
    update_product as gql_update_product,
    delete_product as gql_delete_product,
    restock_product as gql_restock_product,
    set_product_availability as gql_set_product_availability,

    # Orders
    create_order as gql_create_order,
    get_orders_for_user as gql_get_orders_for_user,
    get_order_by_id as gql_get_order_by_id,
    update_order as gql_update_order,
    delete_order as gql_delete_order,
    set_order_status as gql_set_order_status,
    change_order_quantity as gql_change_order_quantity
)


"""
===========================
 USER HANDLERS (Admin + Customer)
===========================
"""

def handle_create_user(name, email, password,  active):
    variables = {
        "input": {
            "name": name,
            "email": email,
            "password": password,
            "role": "CUSTOMER",
            "active": active,
        }
    }
    
    result = gql_create_user(variables["input"])
    
    if result is None:
        return "❌ Error: No response from GraphQL server", "", "", "",  False

    if "errors" in result:
        return f"❌ Error: {result['errors'][0]['message']}", "", "", "",  False

    if "createUser" not in result:
        return f"❌ Error: Unexpected response format: {result}", "", "", "", False

    user = result["createUser"]
    message = (
        f"✅ User Created!\n\n"
        f"ID: {user['id']}\n"
        f"Name: {user['name']}\n"
        f"Email: {user['email']}\n"
        # f"Role: {user['role']}\n"
        f"Active: {user['active']}"
    )

    # Return new message + clear form fields for Gradio
    return message, "", "", "", False

# ADMIN
def handle_list_users(event: gr.SelectData):
    pass
    # event.value → the selected cell value
    # event.index → (row_index, col_index)

    # _, col = event.index

    # # Only allow clicking on the first column (user ID)
    # if col != 0:
    #     return "", "❌ You can only select by clicking the User ID column."

    # user_id = event.value
    # return user_id, f"✅ Selected User ID: {user_id}"

# Admin-only (stub)
def handle_get_user(user_id):
    pass  

# Admin-only (stub)
def handle_update_user(user_id, name, email):
    pass

# Admin-only (stub)
def handle_delete_user(user_id):
    pass


"""
===========================
 PRODUCT HANDLERS (Admin + Customer)
===========================
"""

# Admin-only (stub)
def handle_add_product(name, price, description, inventory):
    pass  

def handle_list_products(after=None, first=10):
    result = gql_get_products_cursor(after, first)

    if result is None or "error" in result:
        return [], None, []

    data = result.get("productsCursor")
    edges = data["edges"]

    rows = [
        [edge["node"]["id"], edge["node"]["name"], edge["node"]["price"]]
        for edge in edges
    ]

    next_cursor = data["pageInfo"]["endCursor"]

    return rows, next_cursor, rows

def handle_get_product(product_id):
    result = gql_get_product_by_id(product_id)
    
    if "error" in result:
        return f"❌ Error: {result['error']}"
    if result.get("product") is None:
        return "❌ Product not found"
    
    p = result["product"]
    message = (
        f"🛍️ Product Details\n\n"
        f"ID: {p['id']}\n"
        f"Name: {p['name']}\n"
        f"Price: ${p['price']}\n"
        f"Description: {p['description']}\n"
        f"Inventory: {p['inventory']}\n"
        f"Available: {p['available']}"
    )

    return message

# Admin-only (stub)
def handle_update_product(product_id, name, price, description, inventory):
    pass  

# Admin-only (stub)
def handle_delete_product(product_id):
    pass  

# Admin-only (stub)
def handle_restock_product(product_id, quantity):
    pass  

# Admin-only (stub)
def handle_set_product_availability(product_id, available):
    pass  



"""
===========================
 ORDER HANDLERS (Admin + Customer)
===========================
"""

def handle_create_order(user_id, product_ids, quantity, total_price, status, created_at):
    variables = {
        "input":{
            "userId": user_id,
            "productIds": product_ids,
            "quantity": quantity,
            "totalPrice": total_price,
            "status": status,
            "createdAt": created_at,
        }
    }
    
    result = gql_create_order(variables["input"])
    
    if result is None:
        return "❌ Error: No response from GraphQL server"
    if "error" in result:
        return f"❌ Error: {result['error']}"
    if "createOrder" not in result:
        return f"❌ Error: Unexpected response format: {result}"
    
    order = result["createOrder"]
    
    message = (
        f"🧾 Order Created!\n\n"
        f"Order ID: {order['id']}\n"
        f"User ID: {order['userId']}\n"
        f"Products: {order['products']}\n"
        f"Quantity: {order['quantity']}\n"
        f"Total Price: {order['totalPrice']}\n"
        f"Status: {order['status']}\n"
        f"Created At: {order['createdAt']}\n"
    )

    return message

# Admin-only (stub)
def handle_list_orders():
    pass  

def handle_get_orders_for_user(user_id):
    result = gql_get_orders_for_user(user_id)
    
    if result is None:
        return "❌ Error: No response from server"
    if "error" in result:
        return f"❌ Error: {result['error']}"
    
    orders = result.get("ordersByUser")
    if orders is None:
        return "❌ Unexpected response format"
    if len(orders)== 0:
        return "This user has no orders"
    
    messages = []
    for o in orders:
        products = ", ".join([p["name"] for p in o["products"]])
        
        messages.append(
            f"🧾 Order {o['id']}\n"
            f"User: {o['userId']}\n"
            f"Products: {products}\n"
            f"Quantity: {o['quantity']}\n"
            f"Total Price: ${o['totalPrice']}\n"
            f"Status: {o['status']}\n"
            f"Created At: {o['createdAt']}\n"
            f"{'-'*40}"
        )

    return "\n".join(messages)

# Admin-only (stub)
def handle_get_order(order_id):
    pass  

# Admin-only (stub)
def handle_update_order(order_id, quantity=None, total_price=None, status=None):
    pass  

# Admin-only (stub)
def handle_delete_order(order_id, user_id):
    pass  

# Admin-only (stub)
def handle_set_order_status(order_id, status):
    pass  

# Admin-only (stub)
def handle_change_order_quantity(order_id, quantity):
    pass  