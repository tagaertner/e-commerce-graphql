import requests
import os

GQL_ENDPOINT = f"http://{os.getenv('GRAPHQL_ENDPOINT', 'gateway:4000')}/query"

# GQL_ENDPOINT = os.getenv("GQL_ENDPOINT", "http://localhost:4000/graphql")


def gql_request(query, variables=None):
    try:
        response = requests.post(
            GQL_ENDPOINT,
            json={"query": query, "variables": variables or {}},
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()
        data = response.json()

        if "errors" in data:
            return {"error": data["errors"][0]["message"]}

        return data["data"]

    except Exception as e:
        return {"error": f"❌ GraphQL request failed: {e}"}

"""
====================================
 USERS
====================================
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
    return gql_request(mutation, {"input": input_data})


# Admin-only (stub)
def get_users():
    pass


# Admin-only (stub)
def get_user_by_id(user_id):
    pass


def update_user(input_data):
    mutation = """
    mutation UpdateUser($input: UpdateUserInput!){
        updatedUser(input: $input){
            id
            name
            email
    
        }
    }
    """
    return gql_request(mutation, {"input": input_data})


def delete_user(user_id):
    mutation = """
    mutation DeleteUser($id: ID!){
        deleteUser(id: $id)
    }
    """
    return gql_request(mutation, {"input": user_id})


"""
====================================
 PRODUCTS 
====================================
"""

# Customer
def get_products_cursor(after=None, first=10):
    query = """
    query GetProductsCursor($after: String, $first: Int) {
        productsCursor(after: $after, first: $first) {
            edges {
                cursor
                node {
                    id
                    name
                    price
                    description
                    inventory
                    available
                }
            }
            pageInfo {
                hasNextPage
                endCursor
            }
            totalCount
        }
    }
    """

    variables = {"after": after, "first": first}
    return gql_request(query, variables)


# Customer
def get_product_by_id(product_id):
    query = """
    query GetProduct($id: ID!){
        product(id: $id){
            id
            name
            price
            description
            inventory
            available
        }
    }
    """
    return gql_request(query, {"id":product_id})


# Admin-only
def create_product(input_data):
    pass


# Admin-only
def update_product(product_id, input_data):
    pass


# Admin-only
def delete_product(input_data):
    pass


# Admin-only
def restock_product(product_id, quantity):
    pass


# Admin-only
def set_product_availability(product_id, available):
    pass


"""
====================================
 ORDERS 
====================================
"""

# Customer
def create_order(input_data):
    mutation = """
    mutation CreateOrder($input: CreateOrderInput!) {
        createOrder(input: $input){
            userId
            quantity
            totalPrice
            status
            createdAt
            products {
                id
                name
            }
        }
    }
    
    """ 
    return gql_request(mutation, {"input": input_data})


# Customer
def get_orders_for_user(user_id):
    query = """
    query GetOrdersForUser($id: ID!) {
        ordersByUser(userId: $id) {
            id
            userId
            quantity
            totalPrice
            status
            createdAt
            products {
                id
                name
            }
        }
    }
    """
    return gql_request(query, {"id": user_id})


# Admin-only
def get_orders():
    pass


# Admin-only
def get_order_by_id(order_id):
    pass


# Admin-only
def update_order(input_data):
    pass


# Admin-only
def delete_order(input_data):
    pass


# Admin-only
def set_order_status(order_id, status):
    pass


# Admin-only
def change_order_quantity(order_id, quantity):
    pass