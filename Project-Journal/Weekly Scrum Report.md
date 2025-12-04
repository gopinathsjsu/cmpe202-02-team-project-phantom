#Backlogs and Burndown Charts
https://docs.google.com/spreadsheets/d/1D12PLO67nH4E2KtAujLxz88DP8UMz1MJgKTrcg9MsaY/edit?usp=sharing

# Weekly Scrum Report  
Project: Marketplace Web Application  
Duration: 6 Sprints (12 Weeks)

---

# WEEK 1 — Sprint 1 (Week 1 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Set up orchestrator service structure  
- Initial backend module wiring  
- MongoDB + basic chat persistence scaffolding  

### **Nikhil**
- Initial WebSocket connection tests on frontend  
- Base chat interface placeholder screens  
- Early navigation skeleton  
- Created initial frontend architecture 

### **Dan**
- Basic page routing + layout  
- Started UI components for listings & homepage  

## ✔ What we plan to work on next
- Completing chat infrastructure  
- Integrating orchestrator → services  
- Building first functional UI pages  

## ✔ Blockers
- None this week — blockers resolved  

---

# WEEK 2 — Sprint 1 (Week 2 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Completed WebSocket server & chat event routing  
- Implemented message persistence in MongoDB  

### **Nikhil**
- Connected frontend chat UI → WebSocket  
- Tested live messaging behavior  

### **Dan**
- Completed initial UI v1.0  
- Integrated homepage shell + listing placeholders  
- Added early price filter components  

## ✔ What we plan to work on next
- User management backend  
- Listings service integration  
- Preparing DB schemas and seed files  

## ✔ Blockers
- None this week — blockers resolved  


---

# WEEK 3 — Sprint 2 (Week 1 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Database schema design for users & listings  
- Listing–User linking via UUID  
- Orchestrator integration cleanup  

### **Nikhil**
- Improved chat UI responsiveness  
- Early AI search component scaffolding  

### **Dan**
- Listing creation UI  
- Initial listing backend integration  
- Seed data folder created  

## ✔ What we plan to work on next
- Implement media upload and blob storage  
- Finish listing CRUD end-to-end  
- Prepare orchestrator → listing routing  

## ✔ Blockers
- Media upload blocked until SAS URL logic finalized
- AI search depends on search-agent endpoint 

---

# WEEK 4 — Sprint 2 (Week 2 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Azure Blob Storage integration  
- Endpoint fixes for listing CRUD  
- UUID-based search improvements  

### **Nikhil**
- Search bar navigation integration  
- WebSocket reconnection improvements  

### **Dan**
- Edit/Delete listing UI  
- Integrated media upload on create-listing  
- Completed seed data & DB integration  

## ✔ What we plan to work on next
- Authentication flow improvements  
- Flagging and reporting system  
- Profile pages  

## ✔ Blockers
- Some media tests blocked on Blob upload timing

---

# WEEK 5 — Sprint 3 (Week 1 of 2)

## ✔ What we worked on / completed
### **Kunal**
- User search API completed  
- Backend authentication cleanup  

### **Nikhil**
- Chat UI finalized  
- Complete message thread UI overhaul  
- Integrated user search into messaging  

### **Dan**
- Login/Signup flow completed  
- Profile page implementation started  

## ✔ What we plan to work on next
- Flagged listing feature  
- Admin dashboard display  
- Enhanced listing filters  

## ✔ Blockers
- Profile editing blocked on backend route availability 

---

# WEEK 6 — Sprint 3 (Week 2 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Finalized user search + improved validations  
- Bug fixes for real-time message delivery  

### **Nikhil**
- Finished messaging feature end-to-end  
- Completed search + AI search UI hooks  

### **Dan**
- Finished profile edit feature  
- Completed ALL user management UI  
- Improved price input filter system  

## ✔ What we plan to work on next
- Reporting/flagging system  
- Admin view for managing content  
- Homepage feature sections  

## ✔ Blockers
- None this week — blockers resolved  

---

# WEEK 7 — Sprint 4 (Week 1 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Backend flagging logic  
- Admin routes for moderation  

### **Nikhil**
- UI for chat search + user selection  
- Optimizations to chat loading speed  

### **Dan**
- Frontend for flagging listings  
- Admin dashboard prototype  
- Personal listings tab  

## ✔ What we plan to work on next
- Completing admin dashboard  
- Improving flagged listings UI  
- Implementing analytics  

## ✔ Blockers
- Admin dashboard backend waiting on flagging API completion   

---

# WEEK 8 — Sprint 4 (Week 2 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Completed moderation APIs  
- Database cleanup scripts  

### **Nikhil**
- UI polish for user search  
- Faster chat refresh improvements  

### **Dan**
- Finished flagged listings UI  
- Completed Admin User Management  
- Bug fixes for listing media display  

## ✔ What we plan to work on next
- Homepage category counts  
- AI search enhancements  
- Saved listings  

## ✔ Blockers
- None  

---

# WEEK 9 — Sprint 5 (Week 1 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Assisted with analytics data piping  
- Backend categorization improvements  

### **Nikhil**
- AI search interface integration  
- Search bar → AI assistant transition logic  

### **Dan**
- Homepage features + category counts  
- Category navigation pages  
- Saved listings backend/JS logic  

## ✔ What we plan to work on next
- Completing analytics dashboard  
- Improving saved listings UX  
- Adding help/about/terms pages  

## ✔ Blockers
- AI search tuning blocked on Q&A dataset  

---

# WEEK 10 — Sprint 5 (Week 2 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Flagged listing seed improvements  
- Analytics query optimizations  

### **Nikhil**
- Finished AI search assistant UI  
- Smooth transitions for search modes  

### **Dan**
- Finished saved listings feature  
- Added native share functionality  
- Completed about/help/terms/privacy pages  

## ✔ What we plan to work on next
- AWS deployment prep  
- Neon Postgres migration  
- MongoDB Atlas migration
- Final stability pass  

## ✔ Blockers
- None  

---

# WEEK 11 — Sprint 6 (Week 1 of 2)

## ✔ What we worked on / completed
### **Kunal**
- Updated listing-service for Neon Postgres SSL
- MongoDB Atlas migration  
- Docker orchestration fixes  

### **Nikhil**
- Final UI improvements across chat + search  
- Cleaned up frontend warnings  

### **Dan**
- Updated Makefile & docker-compose for Neon  
- README deployment documentation  
- Final UX polish  

## ✔ What we plan to work on next
- AWS deployment  
- Load balancer + domain routing setup  
- Environment variable + secrets management  

## ✔ Blockers
- None  

---

# WEEK 12 — Sprint 6 (Week 2 of 2 — Deployment Week)

## ✔ What we worked on / completed
### **Kunal**
- Containerized services for production  
- ECR/ECS/EC2 deployment configuration  

### **Nikhil**
- UI cleanup, performance improvements  
- Final testing of chat + search features  

### **Dan**
- Full AWS deployment  
- Load balancer + routing  
- Final production bug fixes & polish  

## ✔ What we plan to work on next
- Final project handoff  
- Presentation preparation  
- Documentation wrap-up  

## ✔ Blockers
- None — project successfully completed

---

# End of Weekly Scrum Report
