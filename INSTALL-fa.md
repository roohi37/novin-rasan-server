# راهنمای نصب نوین‌رسان (فارسی)

> پیام‌رسان خودمیزبان، بر پایهٔ teamgram-server. این راهنما گام‌به‌گام است.

## پیش‌نیازها
- **Go 1.23+** (برای بیلد)
- **MySQL 8** (دیتابیس اصلی)
- **Redis** (کش و KV)
- **Etcd** (سرویس‌دیسکاوری)
- **MinIO** (ذخیره‌سازی فایل/عکس)
- سیستم‌عامل: Ubuntu 22.04+/Debian (پیشنهادی)

## گام ۱ — دریافت کد
```bash
git clone <repo-url> novin-rasan-server
cd novin-rasan-server
```

## گام ۲ — راه‌اندازی زیرساخت با Docker
```bash
# بالا آوردن MySQL/Redis/Etcd/MinIO
docker compose -f docker-compose-env.yaml up -d
```

## گام ۳ — ساخت دیتابیس
```bash
mysql -h 127.0.0.1 -uroot -p -e "CREATE DATABASE teamgram CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
# ایمپورت اسکیمای اولیه (در teamgramd/deploy/sql)
for f in teamgramd/deploy/sql/*.sql; do mysql -h127.0.0.1 -uroot -p teamgram < "$f"; done
```

## گام ۴ — تنظیم کانفیگ
فایل‌های `teamgramd/etc/*.yaml` را ویرایش کن:
- **MySQL:** `DSN: user:pass@tcp(127.0.0.1:3306)/teamgram?charset=utf8mb4&parseTime=true`
- **Etcd:** `Hosts: [127.0.0.1:2379]`
- **MinIO:** آدرس/کلید/رمز
- **SMS (اختیاری):** در `bff.yaml` → `SmsVerifyCode` (SendCodeUrl / VerifyCodeUrl)
- **پوش آفلاین (اختیاری):** در `session.yaml` و `sync.yaml` → `PushQueueDSN`

## گام ۵ — بیلد
```bash
chmod +x build.sh
./build.sh          # همهٔ سرویس‌ها به teamgramd/bin/
# یا فقط چند سرویس:
make biz sync bff session
```

## گام ۶ — اجرا
```bash
cd teamgramd/bin
./runall2.sh        # اسکریپت اجرای همه
# یا هر سرویس جدا:
./authsession -f=../etc/authsession.yaml &
./biz         -f=../etc/biz.yaml &
# ... (session, sync, msg, bff, gnetway, ...)
```

## گام ۷ — وب‌کلاینت و اپ
- وب‌کلاینت (tweb) را روی nginx سرو کن (پورت 8447 در نمونهٔ ما).
- اپ اندروید: `TMessagesProj` را با Android Studio بیلد کن.

## عیب‌یابی سریع
| مشکل | راه‌حل |
|---|---|
| `panic: FromPeer` | گارد `isUsablePeer` در `app/service/biz/dialog` باید فعال باشد |
| پیام آفلاین نمی‌رسد | `PushQueueDSN` در `session.yaml` + `sync.yaml` ست و جدول `push_queue` ساخته باشد |
| `bad_server_salt` حلقه | فیکس salt-race در session لازم است |
| اتصال اپ برقرار نمی‌شود | پورت‌های 8444/10443/5222 (TCP) و 11443 (WS) باز باشند |

## نکتهٔ امنیتی
- فایل‌های `*.yaml` حاوی رمز دیتابیس/MinIO هستند → هرگز در git commit نکن (`.gitignore` پوشش می‌دهد).
