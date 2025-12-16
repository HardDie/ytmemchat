package alerts

//// WebhookPayload defines the structure of the incoming JSON from the webhook.
//type WebhookPayload struct {
//	Filename string  `json:"filename"` // e.g., "new_follower.gif", "cheer_alert.mp4"
//	Volume   float64 `json:"volume"`   // New field for volume (0.0 to 1.0)
//	Scale    float64 `json:"scale"`    // NEW FIELD for scale multiplier (e.g., 0.5 to 2.0)
//}

//// webhookHandler receives the POST request from the webhook source.
//func webhookHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var payload WebhookPayload
//	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
//		http.Error(w, "Invalid payload", http.StatusBadRequest)
//		return
//	}
//
//	log.Printf("Received webhook data: %s", payload.Filename)
//	broadcast <- payload // Send the payload to the broadcast channel
//
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("Webhook received and broadcasted"))
//}
