import copy
import json

JSON_FILE = 'docs/openapi/swagger.json'

BASE_FILTER = {
    'in': 'query',
    'description': '{{description}}',
    'name': '{{name}}',
    'schema': {
        'type': 'object',
        'properties': {},
    },
    'style': 'deepObject',
}

BASE_FILTER_PROP = {
    'type': '{{type}}',
}


def new_filter(name, props, desc=None):
    filter = copy.deepcopy(BASE_FILTER)
    filter['name'] = name
    if desc: filter['description'] = desc
    if props:
        for p in props:
            filter['schema']['properties'].update(p)
    return filter


def new_filter_prop(name, type):
    prop = copy.deepcopy(BASE_FILTER_PROP)
    prop['type'] = type
    return {
        name: prop
    }


def process_api_filter(api):
    params = api.get('parameters', [])
    if not params:
        return

    new_params = []
    param_filter = []
    param_page = []
    param_sort = None

    for p in params:
        name = p['name']
        if name.startswith('filter['):
            param_filter.append(p)
        elif name.startswith('page['):
            param_page.append(p)
        elif name == 'sort':
            param_sort = p
        else:
            new_params.append(p)

    # Process param_filter/param_page
    if param_filter:
        new_params.append(process_param_filter(param_filter))
    if param_page:
        new_params.append(process_param_page(param_page))
    if param_sort:
        new_params.append(param_sort)

    # Save result
    if new_params:
        api['parameters'] = new_params


def process_param_filter(param_filter):
    props = []
    desc = []
    for filter in param_filter:
        name = filter['name']  # name in form filter[XXX]
        prop_name = name[len('filter['):len(name)-1]
        prop_type = filter['schema']['type']
        desc.append(filter.get('description',''))
        props.append(new_filter_prop(prop_name, prop_type))

    return new_filter(
        name='filter',
        props=props,
        desc=', '.join(desc),
    )


def process_param_page(param_page):
    props = []
    desc = []
    for page in param_page:
        name = page['name']  # name in form page[XXX]
        prop_name = name[len('page['):len(name)-1]
        prop_type = page['schema']['type']
        desc.append(page.get('description',''))
        props.append(new_filter_prop(prop_name, prop_type))

    return new_filter(
        name='page',
        props=props,
        desc=', '.join(desc),
    )


# A path content is a dict of:
# {
#     "get": {},
#     "post": {},
#     "put": {}
# }
def process_path(path, content):
    for method, v in content.items():
        # Process the filter style filter[XXX]=XXX&page[YYY]&sort=ZZZ
        process_api_filter(v)

# Sort all arrays to ensure deterministic output
def sort_array_in_json_recursively(obj):
    if isinstance(obj, dict):
        return {key: sort_array_in_json_recursively(obj[key]) for key in obj.keys()}
    elif isinstance(obj, list):
        processed_list = [sort_array_in_json_recursively(item) for item in obj]
        try:
            return sorted(processed_list)
        except TypeError:
            return processed_list
    else:
        return obj

def main():
    # Load swagger.json file
    with open(JSON_FILE) as f:
        data = json.load(f)

    for k, v in data['paths'].items():
        process_path(k, v)
    data = sort_array_in_json_recursively(data)

    # Save swagger.json file
    with open(JSON_FILE, 'w') as f:
        json.dump(data, f, indent=2)

    print('Support additional features of OpenAPI v3: DONE.')


main()
