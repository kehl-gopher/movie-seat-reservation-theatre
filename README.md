# 🎬 Movy — Movie Seat Reservation System

**Movy** is a full-featured backend system for a cinema to manage and streamline movie scheduling, seat reservations, and ticketing. Designed for **a single theater (admin)**, users can register, browse upcoming shows, reserve seats, pay securely, and receive scannable QR code tickets. Built with **Go**, **PostgreSQL**, **Redis**, and modern authentication/payment systems.

---

## 📦 Tech Stack

| Layer          | Technology                  |
|----------------|-----------------------------|
| Language       | Go (Golang)                 |
| Database       | PostgreSQL                  |
| Caching        | Redis                       |
| Payment        | Paystack                    |

---

## 🛠 Core Features

### 🎥 Movie & Show Management
- Admin can create and manage:
  - Movies
  - Halls with custom seat layouts (rows × seats)
  - Showtimes (with start time, end time, and pricing)

### 👥 User Interaction
- User registration/login via Email or Google
- Email verification required before booking
- View available movies and showtimes

### 💺 Seat Booking Flow
1. Select a movie and preferred showtime
2. Pick available seats
3. Proceed to Paystack payment
4. On success:
   - Booking is confirmed
   - Seats are locked
   - Ticket is issued

### 🧾 Ticketing
- Each ticket includes:
  - Movie name, showtime, hall, seat number(s)
  - A unique QR code
- QR code can be scanned by the cinema to verify authenticity

### 💰 Payments
- Paystack integration for secure online payments
- Transaction verification and confirmation hooks

### 🔒 Security & Caching
- JWT-secured APIs with refresh token handling
- Redis used to:
  - Cache seat availability
  - Prevent double-bookings
  - Improve performance during high demand

---

## ⚙️ How It Works (System Logic)

| Flow               | Description                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| **Auth**           | JWT , email verification, Google OAuth2 login               |
| **Booking**        | Seats held temporarily in Redis, then confirmed post-payment                |
| **QR Code**        | Unique per ticket, scanned at cinema entrance for verification              |
| **Pricing**        | Controlled per show by admin; stored as `Decimal(10,2)` in DB               |
| **Time Validation**| Start/end times validated against current time and show date                |

---

## 🧪 Running Locally

### 1. Clone the repository
```bash
git clone https://github.com/your-username/movy.git
cd movy
```

### 2. Create `.env` configuration
```env
DATABASE_URL=postgres://username:password@localhost:5432/movy_db
REDIS_URL=redis://localhost:6379
JWT_SECRET=supersecretkey
PAYSTACK_SECRET_KEY=your_paystack_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

### 3. Run the app
```bash
go run main.go
```

---

## 📬 Contact
For questions, suggestions, or support, please open an issue or reach out to the project maintainer (kelanidarasimi9@gmail.com).


# NB: This repo is under construction
