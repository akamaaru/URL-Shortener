This is a server application on Go for adding and using aliases (random or chosen) to redirect to a certain web page.

## Stack
- Router and middlewares: `chi/v5` and `chi/v5/middleware`
- Authentication: `BasicAuth`
- Environment variables: `godotenv`
- Database: `modernc.org/sqlite`
- Logger: `log/slog`
- Validation: `validator/v10`
- Test: `testing`, `testify/mock`, `net/http/httptest`

## Requests
- `POST /url {"url": "https://www.link.com/", "alias": "alias"}`
- `GET /alias`
- `DELETE /url/alias`
