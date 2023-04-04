package hhttp

const (
	// MethodGet HTTP method
	GET = "GET"
	// MethodPost HTTP method
	POST = "POST"
	// MethodPut HTTP method
	PUT = "PUT"
	// MethodDelete HTTP method
	DELETE = "DELETE"
	// MethodDelete HTTP method
	PATCH = "PATCH"
)
const (
	SerializationType string = "Content-Type"

	// SerializationTypeFormData 常见在表单的文件上传
	SerializationTypeFormData string = "multipart/form-data"

	SerializationTypeJSON    string = "application/json"
	SerializationTypeWWWFrom string = "application/x-www-form-urlencoded"

	// SerializationTypePlainText 用于请求和响应文本数据。
	SerializationTypePlainText string = "text/plain; charset=utf-8"
	// SerializationTypeImageGIF gif图片
	SerializationTypeImageGIF string = "image/gif"
	// SerializationTypeImageJPEG jpeg图片
	SerializationTypeImageJPEG string = "image/jpeg"
	// SerializationTypeImagePNG png图片
	SerializationTypeImagePNG string = "image/png"
)
