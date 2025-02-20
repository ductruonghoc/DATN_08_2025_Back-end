# DATN_08_2024_Back-end
 Master: Lê Minh Đức

## **👑 Master: Lê Minh Đức**  

This project follows a structured organization to ensure **maintainability, scalability, and readability**. Below is a breakdown of the directory structure and its purpose.  

---

## **📂 Directory Structure**  

### **1️⃣ `./controllers` – Logic Handlers**  
Contains all **handlers** responsible for processing logic.  
- Handlers that belong to the **same group** are placed in the same controller file.  

---

### **2️⃣ `./models` – Database Layer**  
- Establishes **database connections**.  
- Maps **Go structs** to database **tables**.  

---

### **3️⃣ `./routes` – API Routes**  
Defines **routing** and maps requests to their corresponding **handlers**.  
- `route.go` → Defines **global routes**.  
- Other files → Define **routes for specific groups**.  

---

### **4️⃣ `./middlewares` – Middleware Processing**  
Contains **middleware functions** that perform **pre-processing** before requests reach the handlers.  

---

### **5️⃣ `./config` – Configuration Files**  
Stores necessary **configuration files** required for **deployment and setup**.  

---

## **📜 Important Files**  

### **`go.mod` – Module & Dependencies**  
- Declares **project dependencies**.  
- Managed using terminal commands like:  
  ```sh
  go mod tidy

### **`go.sum` - Dependencies checksum (Do not reach)**


### **`main.go` - Only application entry point**