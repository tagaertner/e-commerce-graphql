"""
handlers.py
Outline for transforming UI inputs → GraphQL client calls.
"""

from graphql_client import (
    create_user, get_users, get_user_by_id, update_user, delete_user,
    add_product, get_products, get_product_by_id, update_product, delete_product,
    create_order, get_orders, get_orders_for_user, get_order_by_id, update_order, delete_order
)

"""
===========================
 USER HANDLERS
===========================
"""

def handle_create_user(name, email):
    pass

def handle_list_users():
    pass

def handle_get_user(user_id):
    pass

def handle_update_user(user_id, name, email):
    pass

def handle_delete_user(user_id):
    pass


"""
===========================
 PRODUCT HANDLERS
===========================
"""

def handle_add_product(name, price, description):
    pass

def handle_list_products():
    pass

def handle_get_product(product_id):
    pass

def handle_update_product(product_id, name, price, description):
    pass

def handle_delete_product(product_id):
    pass


"""
===========================
 ORDER HANDLERS
===========================
"""

def handle_create_order(user_id, product_id, quantity):
    pass

def handle_list_orders():
    pass

def handle_list_orders_for_user(user_id):
    pass

def handle_get_order(order_id):
    pass

def handle_update_order(order_id, quantity):
    pass

def handle_delete_order(order_id):
    pass