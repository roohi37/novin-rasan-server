# مشارکت در نوین‌رسان

ممنون که می‌خواهی کمک کنی! 🌱

## قواعد
- زبان کد و کامیت‌ها: انگلیسی یا فارسی (یکدست).
- هر تغییر = یک PR با توضیح روشن.
- تست‌ها/بیلد باید سبز باشند.

## گام‌ها
1. Fork کن.
2. شاخه بساز: `git checkout -b feature/my-change`
3. کامیت کن: `git commit -m "feat: توضیح"`
4. Push کن و PR باز کن.

## نکات مهم
- ⚠️ **هرگز** فایل‌های حاوی رمز/کلید (`*.yaml` با DSN، `.secrets/`) را commit نکن.
- اگر قابلیت جدید اضافه می‌کنی، در `NOTICE` توضیح بده.
- مجوز پروژه Apache-2.0 است؛ تغییرات تو هم زیر همان مجوز منتشر می‌شوند.

## ساختار
- `app/service/` — سرویس‌های پایه (biz, authsession, media, ...)
- `app/messenger/` — msg, sync
- `app/interface/` — session, gnetway
- `app/bff/` — لایهٔ BFF (نزدیک به کلاینت)
- `teamgramd/etc/` — کانفیگ‌ها
- `teamgramd/deploy/sql/` — اسکیمای دیتابیس
