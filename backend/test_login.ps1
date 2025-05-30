$body = @{
    username = "admin"
    password = "admin123"
    platform = "oa"
    device_type = "web"
    device_name = "test"
    device_id = "test123"
    app_version = "1.0.0"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/oa/auth/login" -Method POST -ContentType "application/json" -Body $body
    Write-Host "Status Code: $($response.StatusCode)"
    Write-Host "Response Body: $($response.Content)"
} catch {
    Write-Host "Error: $($_.Exception.Message)"
    Write-Host "Response: $($_.Exception.Response)"
} 