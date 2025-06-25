package com.example.generator.controller;

import com.example.generator.resourceFile.MultipartInputStreamFileResource;
import org.apache.coyote.Response;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.http.*;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.multipart.MultipartFile;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api")
public class ImageGenerationController {

    @Autowired
    private RestTemplate restTemplate;

    @PostMapping("/edit-image")
    public ResponseEntity<byte[]> generateImage(
            @RequestParam("image") MultipartFile image,
            @RequestParam("description") String description) {

        try {
            // Prepare request to Go service
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
            body.add("description", description);
            body.add("image", new MultipartInputStreamFileResource(image.getInputStream(), image.getOriginalFilename()));

            HttpEntity<MultiValueMap<String, Object>> requestEntity = new HttpEntity<>(body, headers);
            RestTemplate restTemplate = new RestTemplate();

            ResponseEntity<byte[]> response = restTemplate.exchange(
                    "http://localhost:9095/edit-image", // this should match Go route
                    HttpMethod.POST,
                    requestEntity,
                    byte[].class
            );

            return ResponseEntity.ok()
                    .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=generated.png")
                    .contentType(MediaType.IMAGE_PNG)
                    .body(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build();
        }
    }

    @PostMapping("/text-to-image")
    public ResponseEntity<byte[]> generateFromText(@RequestParam("description") String description) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);

            MultiValueMap<String, String> map = new LinkedMultiValueMap<>();
            map.add("description", description);

            HttpEntity<MultiValueMap<String, String>> request = new HttpEntity<>(map, headers);

            RestTemplate restTemplate = new RestTemplate();
            ResponseEntity<byte[]> response = restTemplate.exchange(
                    "http://localhost:9095/text-to-image", // this should match Go route
                    HttpMethod.POST,
                    request,
                    byte[].class
            );

            return ResponseEntity.ok()
                    .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=generated.png")
                    .contentType(MediaType.IMAGE_PNG)
                    .body(response.getBody());
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build();
        }
    }

    @PostMapping("/text-to-text")
    public ResponseEntity<String> generateTextOnly(@RequestParam("description") String description) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);

            MultiValueMap<String, String> map = new LinkedMultiValueMap<>();
            map.add("description", description);

            HttpEntity<MultiValueMap<String, String>> request = new HttpEntity<>(map, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/text-to-text", // this should match Go route
                    HttpMethod.POST,
                    request,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Error generating text");
        }
    }

    @PostMapping("/generate-multi-image")
    public ResponseEntity<String> generateMultipleImages(
            @RequestParam("image") MultipartFile image,
            @RequestParam("description") String description) {

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
            body.add("description", description);
            body.add("image", new MultipartInputStreamFileResource(image.getInputStream(), image.getOriginalFilename()));

            HttpEntity<MultiValueMap<String, Object>> requestEntity = new HttpEntity<>(body, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/generate-multi-image",  // this should match Go route
                    HttpMethod.POST,
                    requestEntity,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("{\"error\":\"Image generation failed\"}");
        }
    }

    @PostMapping("/code-execute")
    public ResponseEntity<String> executePromptCode(@RequestParam("prompt") String prompt) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);

            MultiValueMap<String, String> formData = new LinkedMultiValueMap<>();
            formData.add("prompt", prompt);

            HttpEntity<MultiValueMap<String, String>> request = new HttpEntity<>(formData, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/code-execute",
                    HttpMethod.POST,
                    request,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("{\"error\":\"Code execution failed\"}");
        }
    }

    @PostMapping("/image-understanding-by-url")
    public ResponseEntity<String> summarizeImageWithPrompt(
            @RequestParam("imageUrl") String imageUrl,
            @RequestParam("prompt") String prompt) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);

            MultiValueMap<String, String> form = new LinkedMultiValueMap<>();
            form.add("imageUrl", imageUrl);
            form.add("prompt", prompt);

            HttpEntity<MultiValueMap<String, String>> request = new HttpEntity<>(form, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/image-understanding-by-url",
                    HttpMethod.POST,
                    request,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("{\"error\":\"Image summarization failed\"}");
        }
    }


    @PostMapping("/summarize-youtubeVideo-from-url")
    public ResponseEntity<String> summarizeVideoWithPrompt(
            @RequestParam("videoUrl") String videoUrl,
            @RequestParam("prompt") String prompt) {

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);

            MultiValueMap<String, String> form = new LinkedMultiValueMap<>();
            form.add("videoUrl", videoUrl);
            form.add("prompt", prompt);

            HttpEntity<MultiValueMap<String, String>> requestEntity = new HttpEntity<>(form, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/ summarize-youtubeVideo-from-url",
                    HttpMethod.POST,
                    requestEntity,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("{\"error\":\"Video summarization failed\"}");
        }
    }

    // only understand PDF
    @PostMapping("/summarize-pdf-from-locally-stored")
    public ResponseEntity<String> summarizePdf(
            @RequestParam("pdf") MultipartFile pdf,
            @RequestParam(value = "description", required = false, defaultValue = "Give me a summary of this PDF file.") String description) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
            body.add("description", description);
            body.add("pdf", new MultipartInputStreamFileResource(pdf.getInputStream(), pdf.getOriginalFilename()));

            HttpEntity<MultiValueMap<String, Object>> requestEntity = new HttpEntity<>(body, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/summarize-pdf-from-locally-stored",
                    HttpMethod.POST,
                    requestEntity,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error summarizing PDF");
        }
    }

    @PostMapping("/summarize-pdf-from-url")
    public ResponseEntity<String> summarizePdfFromUrl(@RequestParam("url") String url,
                                                      @RequestParam(defaultValue = "Summarize this document") String prompt) {
        try {
            // Download PDF from the given URL
            RestTemplate restTemplate = new RestTemplate();
            ResponseEntity<byte[]> pdfResponse = restTemplate.getForEntity(url, byte[].class);

            if (pdfResponse.getStatusCode() != HttpStatus.OK || pdfResponse.getBody() == null) {
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("{\"error\":\"Failed to download PDF\"}");
            }

            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            ByteArrayResource pdfResource = new ByteArrayResource(pdfResponse.getBody()) {
                @Override
                public String getFilename() {
                    return "file.pdf";
                }
            };

            MultiValueMap<String, Object> formData = new LinkedMultiValueMap<>();
            formData.add("prompt", prompt);
            formData.add("file", pdfResource);

            HttpEntity<MultiValueMap<String, Object>> request = new HttpEntity<>(formData, headers);

            ResponseEntity<String> response = restTemplate.exchange(
                    "http://localhost:9095/summarize-pdf-from-url",
                    HttpMethod.POST,
                    request,
                    String.class
            );

            return ResponseEntity.ok(response.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("{\"error\":\"Failed to summarize PDF\"}");
        }
    }

    @PostMapping("/summarize-image-from-locally-stored")
    public ResponseEntity<String> captionImage(
            @RequestParam("image") MultipartFile imageFile,
            @RequestParam("prompt") String prompt
    ) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
            body.add("prompt", prompt);
            body.add("image", new MultipartInputStreamFileResource(imageFile.getInputStream(), imageFile.getOriginalFilename()));

            HttpEntity<MultiValueMap<String, Object>> requestEntity = new HttpEntity<>(body, headers);

            RestTemplate restTemplate = new RestTemplate();
            String goApiUrl = "http://localhost:9095/summarize-image-from-locally";

            ResponseEntity<String> response = restTemplate.postForEntity(goApiUrl, requestEntity, String.class);
            return ResponseEntity.ok(response.getBody());
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Error: " + e.getMessage());
        }
    }

    @PostMapping("/compare-two-images")
    public ResponseEntity<String> compareImages(
            @RequestParam("image1") MultipartFile image1,
            @RequestParam("image2") MultipartFile image2,
            @RequestParam("prompt") String prompt
    ) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
            body.add("prompt", prompt);
            body.add("image1", new MultipartInputStreamFileResource(image1.getInputStream(), image1.getOriginalFilename()));
            body.add("image2", new MultipartInputStreamFileResource(image2.getInputStream(), image2.getOriginalFilename()));

            HttpEntity<MultiValueMap<String, Object>> requestEntity = new HttpEntity<>(body, headers);

            RestTemplate restTemplate = new RestTemplate();
            String url = "http://localhost:9095/compare-two-images"; // this should match Go route

            ResponseEntity<String> response = restTemplate.postForEntity(url, requestEntity, String.class);
            return ResponseEntity.ok(response.getBody());
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Error: " + e.getMessage());
        }
    }

}
