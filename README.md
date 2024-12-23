# *Movie Booking Backend Application*

This is a backend application designed for seamless movie ticket booking and management. The system is built with a microservices architecture, comprising six services that handle different functionalities to provide an efficient and scalable user experience.

---

## *Features*
- Movie listing and management.
- Real-time seat selection and availability updates.
- Secure payment processing.
- User authentication and profile management.
- Comprehensive booking history and ticket generation.
- Scalable architecture with six independent services.

---

## *Microservices Overview*

### *1. User/Admin Service* -->  [repo](https://github.com/aparnasukesh/bookyourshow-api-gateway-service.git)
- Manages user registration, login, and profile data.

### *2. Movie-Booking Service* -->  [repo](https://github.com/aparnasukesh/bookyourshow-movies-booking-service.git)
- Stores and manages movie details like title, genre, duration, and showtimes.
- Theater management
- It manages the relationship between theaters, movies, showtimes and bookings.
- Provides APIs for listing and searching movies, theaters, showtimes.

### *4. Payment Service* -->  [repo](https://github.com/aparnasukesh/bookyourshow-payment-service.git)
- Processes payments securely using integrated payment gateway - Razorpay.
- Handles payment statuses and refunds.

### *5. Notification Service* -->  [repo](https://github.com/aparnasukesh/bookyourshow-notification-service.git)
- Sends booking confirmations and updates via email.
- Manages notification preferences for users.
- Help-desk chat option for users by implementing websocket and rabbitmq.

### *6. Auth Service* -->  [repo](https://github.com/aparnasukesh/bookyourshow-authentication-service.git)
- Enables theater owners and admins to manage movies, theaters, schedules, and seat layouts.
- Handles authentication and authorization.

### *1. Api-Getway Service* -->  [repo](https://github.com/aparnasukesh/bookyourshow-api-gateway-service.git)
- The API Gateway Service routes client requests to the correct backend services, acting as the central entry point.  
- It ensures secure communication, handles authentication, and provides a unified interface for all microservices.  


---

## *Tech Stack*
- *Programming Language*: Golang  
- *Database*: SQL and MongoDB  
- *Version Control*: Git  
- *Architecture*: Microservices  
- *API Protocol*: REST

---

## *How to Use*
1. Register a user via the *User Service API*.  
2. Log in as a user using the *Auth Service API* for authentication.  
3. Browse available movies using the *Movie Service API*.  
4. Book tickets through the *Booking Service*.  
5. Process payments via the *Payment Service*.  
6. Receive booking notifications from the *Notification Service*.  
7. Admins can manage movies and schedules through the *Admin Service*.  
