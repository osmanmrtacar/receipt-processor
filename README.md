# Receipt Processor

This is API processes receipts and calculates points based on specific rules. The service provides two main endpoints: one for processing receipts and another for retrieving points for a processed receipt.

## API Endpoints

### 1. Process Receipts

- **Path**: `/receipts/process`
- **Method**: `POST`
- **Payload**: Receipt JSON [Example](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/examples/simple-receipt.json)
- **Response**: JSON containing an id for the receipt

This endpoint takes a JSON receipt and returns a JSON object with a generated ID for the receipt.

Example Response:
```json
{
  "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
}
```

### 2. Get Points

- **Path**: `/receipts/{id}/points`
- **Method**: `GET`
- **Response**: A JSON object containing the number of points awarded

This endpoint retrieves the points awarded for a specific receipt ID.

Example Response:
```json
{
  "points": 32
}
```

## Running the Application
```bash
go mod tidy

go run cmd/server/main.go
```
or

```bash
docker build -t receipt-processor .
docker run -p 8080:8080 receipt-processor
```

## Implementation Details

- The application is built using go.
- Data is stored in sqlite.

## Testing
```bash
go test ./...
```
