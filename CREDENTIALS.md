# 🔐 Comply360 - Official Credentials Documentation

**Last Updated:** December 26, 2025  
**Status:** ✅ **Single Source of Truth**

---

## 🎯 **OFFICIAL DEFAULT CREDENTIALS**

### **Application Login (Comply360 Platform)**
```
Email:    admin@comply360.com
Password: Admin@123
Tenant:   9ac5aa3e-91cd-451f-b182-563b0d751dc7
```

**⚠️ IMPORTANT:** 
- Password uses `@` symbol, NOT `!`
- Password is: `Admin@123` (capital A, lowercase dmin, @ symbol, 123)
- This is the ONLY official password for the admin user

---

## 📊 **Password Hash**

The password `Admin@123` should have this bcrypt hash:
```
$2b$12$A6vMH.18MDpMQjeF97TS5OejovEbFiU7QJEvWqakbvbXACRgoUTOG
```

**Hash Details:**
- Algorithm: bcrypt
- Cost: 12
- Format: $2b$12$...

---

## 🗄️ **Database Credentials**

### **PostgreSQL (Main Database)**
```
Host:     localhost
Port:     5432
Database: comply360_dev
User:     comply360
Password: dev_password
```

### **PostgreSQL (Odoo Database)**
```
Host:     localhost
Port:     5432
Database: comply360_odoo
User:     comply360
Password: dev_password
```

---

## 🔧 **Service Credentials**

### **Redis**
```
URL:      redis://localhost:6379/0
Password: (none)
```

### **RabbitMQ**
```
URL:      amqp://comply360:dev_password@localhost:5672/
User:     comply360
Password: dev_password
```

### **MinIO (S3-compatible)**
```
Endpoint:  localhost:9000
Access Key: comply360
Secret Key: dev_password
Bucket:     comply360-documents
```

### **Odoo**
```
URL:      http://localhost:8069
Database: comply360_odoo
Username: admin
Password: admin
```

---

## 📝 **Where This Password Is Used**

### **Scripts Using `Admin@123`:**
1. ✅ `scripts/fix-and-start.sh` - Line 85, 94
2. ✅ `scripts/test-login.sh` - Line 48, 72
3. ✅ All documentation files

### **Scripts Using WRONG Password:**
1. ❌ `scripts/restart-all.sh` - Line 229 uses `Admin123!` (WRONG)

---

## 🔄 **Password Consistency Rules**

1. **ALWAYS use:** `Admin@123` (with `@` symbol)
2. **NEVER use:** `Admin123!` (with `!` symbol)
3. **NEVER use:** `Admin123` (no special character)
4. **NEVER use:** `admin@123` (lowercase)

---

## 🛠️ **How to Verify/Update Password**

### **Check Current Password Hash:**
```sql
SELECT email, LEFT(password_hash, 30) as hash_preview 
FROM tenant_9ac5aa3e91cd451fb182563b0d751dc7.users 
WHERE email = 'admin@comply360.com';
```

### **Update Password Hash (if needed):**
```sql
UPDATE tenant_9ac5aa3e91cd451fb182563b0d751dc7.users 
SET password_hash = '$2b$12$A6vMH.18MDpMQjeF97TS5OejovEbFiU7QJEvWqakbvbXACRgoUTOG',
    failed_login_attempts = 0,
    locked_until = NULL
WHERE email = 'admin@comply360.com';
```

### **Generate New Hash (if needed):**
```python
import bcrypt
hash = bcrypt.hashpw("Admin@123".encode('utf-8'), bcrypt.gensalt(12))
print(hash.decode('utf-8'))
```

---

## ⚠️ **Common Issues**

### **Issue 1: Wrong Password Symbol**
- ❌ Wrong: `Admin123!` (uses `!`)
- ✅ Correct: `Admin@123` (uses `@`)

### **Issue 2: Account Locked**
If account is locked due to failed attempts:
```sql
UPDATE tenant_9ac5aa3e91cd451fb182563b0d751dc7.users 
SET failed_login_attempts = 0, locked_until = NULL 
WHERE email = 'admin@comply360.com';
```

### **Issue 3: Password Hash Mismatch**
If password doesn't work, verify hash matches:
```sql
-- Should return the hash starting with $2b$12$A6vMH...
SELECT password_hash FROM tenant_9ac5aa3e91cd451fb182563b0d751dc7.users 
WHERE email = 'admin@comply360.com';
```

---

## 📚 **Documentation References**

All scripts and documentation should reference:
- **Password:** `Admin@123`
- **Email:** `admin@comply360.com`
- **Tenant ID:** `9ac5aa3e-91cd-451f-b182-563b0d751dc7`

---

## ✅ **Verification Checklist**

- [ ] Database password hash matches `$2b$12$A6vMH...`
- [ ] All scripts use `Admin@123` (not `Admin123!`)
- [ ] Account is not locked (`locked_until` is NULL)
- [ ] Failed login attempts = 0
- [ ] User exists in both `global_users` and tenant schema
- [ ] User has `tenant_admin` role assigned in `user_roles` table

---

**This is the SINGLE SOURCE OF TRUTH for all credentials.**

