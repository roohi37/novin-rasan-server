# نوین‌رسان (Novin Rasan) — پیام‌رسان خودمیزبان

> پیام‌رسان متن‌باز و خودمیزبان با بومی‌سازی کامل فارسی/RTL، بر پایهٔ [teamgram-server](https://github.com/teamgram/teamgram-server) + [tweb](https://github.com/Ajaxy/telegram-tt).
> این نسخه شامل پچ‌های اختصاصی: **پوش آفلاین end-to-end**، فیکس salt-race، و گارد FromPeer.

## ✨ ویژگی‌ها
- **پوش آفلاین end-to-end**: پیام‌های کاربر آفلاین در صف ذخیره و روی reconnect تحویل داده می‌شوند.
- **بومی‌سازی فارسی/RTL کامل** + برندینگ مستقل.
- **کلاینت وب + اپ اندروید** در کنار سرور.
- **باگ‌فیکس‌های تولیدی:** salt-race (session)، گارد FromPeer (dialog)، appcode (bff).

## 🏗️ معماری
```
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│  Web (tweb) │   │  Android    │   │  (iOS)      │
└──────┬──────┘   └──────┬──────┘   └──────┬──────┘
       │  MTProto         │                  │
       └────────┬─────────┴──────────────────┘
                ▼
        ┌───────────────┐
        │   gnetway     │  (TCP/WS ingress)
        └───────┬───────┘
   ┌────────┬───┴────┬────────┬────────┐
   ▼        ▼        ▼        ▼        ▼
 session   sync     msg      bff      biz
   └────────┴────────┴────────┴────────┘
                │
        MySQL + Redis + Etcd + MinIO
```

## 🚀 راه‌اندازی سریع
```bash
# پیش‌نیازها: Go 1.23+, MySQL, Redis, Etcd, MinIO
mysql -e "CREATE DATABASE teamgram CHARACTER SET utf8mb4;"
# کانفیگ: teamgramd/etc/*.yaml (DSN، Etcd، MinIO)
./build.sh
# اجرا با systemd یا docker-compose
```

## 🤝 حمایت مالی (Support)
اگر می‌خواهی از این پروژه حمایت کنی:

- **USDT (شبکه TRC20):** `TEHAA4SmirDXx1F2v2pFKwzZ4PuhQb2WNm`
- **TON:** `UQCUTaeHfGmyinhMhN3t7ByT_o_xxrvEqPFSVvc5ijcixn_x`
- **داخل ایران (ریال):** درگاه پرداخت — [به‌زودی]
- **نسخهٔ Enterprise (پشتیبانی اختصاصی + قابلیت‌های افزوده):** برای دریافت، از طریق Issues/ایمیل تماس بگیر.

> ⚠️ فقط با شبکهٔ اعلام‌شده واریز کنید. واریز در شبکهٔ اشتباه باعث از‌دست‌رفتن دارایی می‌شود.

## 📄 مجوز
- سرور: **Apache License 2.0** (بر پایهٔ teamgram-server)
- کلاینت وب: **GPL-3.0** (بر پایهٔ tweb)
- اپ اندروید: **GPL-2.0** (بر پایهٔ Telegram Android)

اعلان‌های حق نشر در فایل [NOTICE](NOTICE) آمده است.

## 🙏 قدردانی
- [teamgram-server](https://github.com/teamgram/teamgram-server)
- [telegram-tt (tweb)](https://github.com/Ajaxy/telegram-tt)
- [Telegram Android](https://github.com/DrKLO/Telegram)

## ⚠️ سلب مسئولیت
این پروژه یک محصول مستقل است و هیچ ارتباط رسمی با Telegram Messenger Inc. ندارد.
