# AgroSeed REST API

AgroSeed adalah aplikasi backend berbasis Golang yang dikembangkan untuk mendukung sistem manajemen bibit tanaman secara terstruktur dan terintegrasi. Aplikasi ini menyediakan layanan REST API yang memungkinkan pengguna melakukan pengelolaan data bibit, pengaturan stok masuk dan keluar, serta pembuatan laporan stok secara otomatis.

Selain manajemen data, AgroSeed juga memiliki fitur rekomendasi bibit yang membantu menentukan bibit tanaman yang sesuai berdasarkan kondisi lahan, seperti jenis tanah, curah hujan, dan luas lahan. Seluruh proses bisnis dirancang agar mudah diakses melalui endpoint API dengan format pertukaran data JSON.

Sistem ini menggunakan PostgreSQL sebagai basis data utama, Gorilla Mux sebagai router HTTP, serta JWT (JSON Web Token) untuk autentikasi dan otorisasi pengguna. Di dalam pengembangannya, AgroSeed menerapkan konsep Pemrograman Fungsional melalui penggunaan fungsi map, filter, dan reduce pada pengolahan data, sehingga kode menjadi lebih modular, reusable, dan mudah dipelihara.

AgroSeed dikembangkan sebagai bagian dari Project Akhir Mata Kuliah Pemrograman Fungsional dan ditujukan sebagai contoh implementasi REST API modern yang aman, terstruktur, dan relevan dengan kebutuhan sistem informasi pertanian.
