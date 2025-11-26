"""
handlers.py
Transform UI inputs → GraphQL client calls.
"""

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

def create_user(name, email):
    pass

def list_users():
    pass  # ADMIN

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