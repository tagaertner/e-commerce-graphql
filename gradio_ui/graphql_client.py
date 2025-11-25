import requests
import os

GQL_ENDPOINT = os.getenv("GRAPHQL_ENDPOINT", "http://gateway:4000/graphql")

"""
graphql_client.py
Outline for e-commerce GraphQL operations.
"""

"""
===========================
 USERS
===========================
"""

def create_user(input_data):
    mutation = """
    mutation CreateUser($input: CreateUserInput!){
        createUser(input: $input){
            id
            name
            email
            role
            active
        }
    }
    """
    try:
        response = requests.post(
            GQL_ENDPOINT,
            json={"query": mutation, "variables": {"input": input_data}},
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()
        data = response.json()

        if "errors" in data:
            return {"error": data["errors"][0]["message"]}

        return data["data"]

    except Exception as e:
        return {"error": f"❌ Failed to submit user: {e}"}


def get_users():
    query = """
    query {
        users {
            id
            name
            email
        }
    }
    """
    pass


def get_user_by_id(user_id):
    pass

def update_user(input_data):
    pass

def delete_user(user_id):
    pass


"""
===========================
 PRODUCTS
===========================
"""

def add_product(input_data):
    pass

def get_products():
    pass

def get_product_by_id(product_id):
    pass

def update_product(input_data):
    pass

def delete_product(product_id):
    pass


"""
===========================
 ORDERS
===========================
"""

def create_order(input_data):
    pass

def get_orders():
    pass

def get_orders_for_user(user_id):
    pass

def get_order_by_id(order_id):
    pass

def update_order(input_data):
    pass

def delete_order(order_id):
    pass