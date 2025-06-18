
# Spring Boot + Go Gemini Image Generation Service

This is a full-stack image generation project integrating **Spring Boot** (Java) and **Go**.  
It allows users to submit an image and a text prompt to generate a new image using **Gemini Image Generation API** (Google GenAI).

---

## 🛠️ Tech Stack

- **Backend 1 (Go)**: Handles image and text prompt, sends request to Gemini API.
- **Backend 2 (Spring Boot - Java)**: REST API that receives frontend input and forwards it to Go service.
- **Google GenAI Gemini API**: Used for text-to-image generation based on prompt and uploaded image.
- **RestTemplate**: For communication between Spring Boot and Go service.
- **Postman**: For testing API requests.

---

## 📁 Project Structure

```
spring-go-image-generator/
│
├── go-service/
│   ├── .env                      # GEMINI_API_KEY stored here
│   ├── go.mod
│   ├── go.sum
│   └── main.go                   # Go server handling Gemini API requests
│
└── springboot-app/
    ├── pom.xml
    └── src/
        └── main/
            ├── java/
            │   └── com/example/generator/
            │       ├── config/
            │       │   └── RestConfig.java                # RestTemplate Bean
            │       ├── controller/
            │       │   └── ImageGenerationController.java # API Endpoint
            │       ├── resourceFile/
            │       │   └── MultipartInputStreamFileResource.java
            │       └── SpringbootAppApplication.java
            └── resources/
                └── application.properties                 # Server port & upload limits
```

---

## ⚙️ Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/jayesh-Somwanshi/GoogleAIStudio_Edit_Image.git
cd spring-go-image-generator
```

### 2. Set Up Go Service

```bash
cd go-service
go mod tidy
```

Create a `.env` file:

```env
GEMINI_API_KEY=your_gemini_api_key_here
```

Run the Go server:

```bash
go to go project file path: go run main.go
```

It will run on `http://localhost:8081`.

---

### 3. Set Up Spring Boot Application

```bash
cd springboot-app
./mvnw spring-boot:run
```

It will run on `http://localhost:8080`.

---

## 🚀 API Usage

### Endpoint

```
POST /api/image/generate
```

### Request

- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `image`: File (e.g., `.png`, `.jpg`)
  - `description`: Text prompt for branding or image transformation

### Example `curl` Command:

```bash
curl --location 'http://localhost:8080/api/image/generate' --form 'image=@"postman-cloud:///1f04c292-588a-4e50-b067-4cec6d34ef19"' --form 'description="This is a high-resolution image of a perfume bottle. Please add elegant and stylish branding to the bottle. Write the brand name Essence Noir in the center of the bottle in a luxurious serif font (e.g., Didot or similar), in gold or black color, depending on what contrasts best with the bottle's light background. Keep the lighting, reflections, and rest of the bottle untouched for a realistic advertisement look"'
```

### Response

- Returns a `200 OK` with generated `image/png` in binary format (automatically downloaded by browsers).

---

## 📷 Sample Image (Input)

> You can use the sample image provided in `postman-cloud:///...` or upload your own `.png`.

![Sample Input](./go-service/sample.png)

---

## 🧠 Notes

- The Gemini API model used: `gemini-2.0-flash-preview-image-generation`
- Supports high-quality branding overlays and realistic advertisement rendering
- Ensure the image uploaded is in `.png` or `.jpeg` format

---

## 🤝 Contributing

1. Fork the repo
2. Create a new branch
3. Commit changes
4. Submit a pull request

---

