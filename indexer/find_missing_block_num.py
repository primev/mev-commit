from elasticsearch import Elasticsearch, helpers

# Initialize Elasticsearch client with authentication
es = Elasticsearch(
    ["http://localhost:9200"],  # Replace with your Elasticsearch host if different
    basic_auth=("elastic", "mev-commit")
)

# Function to get all numbers using scroll API
def get_all_numbers():
    numbers = []
    scroll_size = 10000
    
    # Initial search request
    response = es.search(
        index="blocks",
        body={
            "size": scroll_size,
            "_source": ["number"],
            "sort": [{"number": "asc"}]
        },
        scroll='2m'
    )
    
    # Get the scroll ID
    scroll_id = response['_scroll_id']
    
    # Get the first batch of numbers
    numbers.extend([hit['_source']['number'] for hit in response['hits']['hits']])
    
    # Continue scrolling until no more hits
    while len(response['hits']['hits']):
        response = es.scroll(scroll_id=scroll_id, scroll='2m')
        numbers.extend([hit['_source']['number'] for hit in response['hits']['hits']])
    
    return numbers

# Get all numbers
all_numbers = get_all_numbers()

# Find missing numbers
missing_numbers = []
for i in range(len(all_numbers) - 1):
    current_number = all_numbers[i]
    next_number = all_numbers[i + 1]
    if next_number != current_number + 1:
        missing_numbers.extend(range(current_number + 1, next_number))

print("Missing numbers:", missing_numbers)
