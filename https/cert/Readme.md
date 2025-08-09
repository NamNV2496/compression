# How to generate cerification for HTTPS

```bash
cd cert

./gen.sh
```

- `client-ext.cnf` is certificate extensions. it will append into final client's certification.
- `server-ext.cnf` is certificate extensions. it will append into final server's certification.

# List generated file


## Certificate Authority files

- `ca-key.pem`: Private key của CA. Giữ bí mật tuyệt đối, dùng để ký chứng chỉ cho server và client.
- `ca-cert.pem`: Public certificate của CA (self-signed). Server và client sẽ dùng file này để xác thực lẫn nhau.

## Client files

- `client-key.pem`: Private key của client. Giữ bí mật, dùng để chứng minh danh tính khi kết nối tới server qua mTLS.
- `client-req.pem`: CSR của client. Tương tự server-req.pem nhưng cho client.
- `client-cert.pem`: Public certificate của client, đã được CA ký và kèm extensions từ client-ext.cnf. Server sẽ kiểm tra chứng chỉ này khi client kết nối.

## Server files
- `server-key.pem`: Private key của server. Dùng để giải mã dữ liệu TLS và tạo chữ ký số khi bắt tay SSL/TLS. Phải bảo mật.
- `server-req.pem`: Certificate Signing Request (CSR) của server. Chứa thông tin public key và subject (CN, OU...) để gửi cho CA ký.
- `server-cert.pem`: Public certificate của server, đã được CA ký và kèm extensions từ server-ext.cnf. Client sẽ kiểm tra chứng chỉ này khi kết nối.

## file khác

- `ca-cert.srl:`: Serial number cuối cùng mà CA đã dùng khi ký chứng chỉ. OpenSSL lưu để đảm bảo mỗi cert được cấp số serial duy nhất.


# Flow

1. CA tự tạo ca-key.pem và ca-cert.pem.
2. Server tạo key riêng và CSR → gửi cho CA → CA ký tạo server-cert.pem.
3. Client tạo key riêng và CSR → gửi cho CA → CA ký tạo client-cert.pem.

## Khi chạy mTLS:

1. Server gửi server-cert.pem cho client.
2. Client gửi client-cert.pem cho server.
3. Cả hai dùng ca-cert.pem để kiểm tra chứng chỉ đối phương.

