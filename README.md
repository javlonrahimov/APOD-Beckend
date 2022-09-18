# APOD-Beckend
### Health-check
| Mathod      | URL pattern | Action |
| ----------- | ----------- | -------|
| GET      | /v1/health-check | Show application information|

### Apods
| Mathod      | URL pattern | Action |
| ----------- | ----------- | -------|
| POST   | /v1/apods | Create a new apod |
| GET   | /v1/apods | Get list of apods |
| GET   | /v1/apods/:id | Show the details of a specific apod |
| PATCH  | /v1/apods/:id | Update specific apod |
| DELETE | /v1/apods/:id| Delete specific apod |

### Users
| Mathod      | URL pattern | Action |
| ----------- | ----------- | -------|
| POST | /v1/users | Register user |
