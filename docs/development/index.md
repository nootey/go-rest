## Development

### Testing API calls with Postman

The api is protected via JWT authentication. While you can just disable the authentication in development,
Postman provides and option to auto login and fetch the bearer token with each request, in the collection.

Here are the steps:

### Setting up a collection

You can do this on a per-request basis, but it's easier to wrap all requests into a collection, so that they use the same
parent mechanism.

- Create a new collection
- Under `Authorization`, set Auth Type to `Bearer Token` and fill the Token field with: `{{jwt_token}}`
- Under variables, creates a new variable `jwt_token`. This is where the token will be held, under current value.
- Under scripts->Pre-request, paste this code:

```js
pm.sendRequest({
    url: "http://localhost:2000/api/v1/auth/login",
    method: "POST",
    header: {
        "Content-Type": "application/json"
    },
    body: {
        mode: "raw",
        raw: JSON.stringify({
            email: "<email>",
            password: "<password>",
        })
    }
}, function (err, res) {
    if (err) {
        console.log("Login request error:", err);
    } else {
        try {
            var jsonResponse = res.json();

            if (res.code === 200 && jsonResponse.token) {
                pm.collectionVariables.set("jwt_token", jsonResponse.token);
                console.log("JWT Token set successfully:", jsonResponse.token);
            } else {
                console.log("Unexpected response structure or missing token:", jsonResponse);
            }
        } catch (error) {
            console.log("Error parsing JSON response:", error);
        }
    }
});
```