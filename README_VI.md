![Go Clean Template](docs/img/logo.svg)

# Go Clean template

[🇬🇧 English](README.md)
[🇨🇳 中文](README_CN.md)
[🇷🇺 RU](README_RU.md)

Mẫu thiết kế Clean Architecture dành cho các dịch vụ Golang

[![Release](https://img.shields.io/github/v/release/minhhoccode111/go-clean-template-gin.svg)](https://github.com/minhhoccode111/go-clean-template-gin/releases/)
[![License](https://img.shields.io/badge/License-MIT-success)](https://github.com/minhhoccode111/go-clean-template-gin/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/minhhoccode111/go-clean-template-gin)](https://goreportcard.com/report/github.com/minhhoccode111/go-clean-template-gin)
[![codecov](https://codecov.io/gh/minhhoccode111/go-clean-template-gin/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/minhhoccode111/go-clean-template-gin)

[![Web Framework](https://img.shields.io/badge/Gin-Web%20Framework-blue)](https://github.com/gin-gonic/gin)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![JSON Handling](https://img.shields.io/badge/Go--JSON-Fast%20Serialization-blue)](https://github.com/goccy/go-json)
[![Query Builder](https://img.shields.io/badge/sqlc-SQL%20Compiler-blue)](https://sqlc.dev/)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/zsais/go-gin-prometheus)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## Tổng quan

Đây là một bản phân nhánh (fork) từ [go-clean-template](https://github.com/evrone/go-clean-template) với các thay đổi:

- Thay thế Fiber bằng Gin
- Thay thế Squirrel bằng Sqlc
- Thêm validatorx wrapper vào `pkg`
- Thêm bộ nhớ đệm (cache) Otter vào `pkg`

Mục đích của dự án mẫu này là để trình diễn:

- Cách tổ chức một dự án để ngăn cấu trúc mã nguồn trở nên chắp vá (spaghetti code).
- Nơi lưu trữ logic nghiệp vụ (business logic) để đảm bảo tính độc lập, rõ ràng và có khả năng mở rộng.
- Cách duy trì quyền kiểm soát khi một vi dịch vụ (microservice) phát triển quy mô lớn.

Dự án sử dụng các nguyên tắc của Robert Martin (được biết đến với tên Uncle Bob).

[Go-clean-template](https://minhhoccode111.com/go-clean-template-gin?utm_source=github&utm_campaign=go-clean-template-gin) được tạo ra và
hỗ trợ bởi [minhhoccode111](https://minhhoccode111.com/?utm_source=github&utm_campaign=go-clean-template-gin).

Mẫu thiết kế này triển khai ba loại máy chủ:

- AMQP RPC (dựa trên RabbitMQ làm [transport](https://github.com/rabbitmq/amqp091-go)
  và [Request-Reply pattern](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))
- MQ RPC (dựa trên NATS làm [transport](https://github.com/nats-io/nats.go)
  và [Request-Reply pattern](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))
- gRPC (framework [gRPC](https://grpc.io/) dựa trên protobuf)
- REST API (framework [Gin](https://github.com/gin-gonic/gin))

## Nội dung

- [Bắt đầu nhanh](#bắt-đầu-nhanh)
- [Cấu trúc dự án](#cấu-trúc-dự-án)
- [Tiêm phụ thuộc (Dependency Injection)](#tiêm-phụ-thuộc-dependency-injection)
- [Kiến trúc Sạch (Clean Architecture)](#kiến-trúc-sạch-clean-architecture)

## Dọn dẹp mẫu (Chỉ dành cho Gin REST API)

Nếu bạn chỉ cần một REST API đơn giản và muốn loại bỏ gRPC, NATS, và RabbitMQ ra khỏi dự án, bạn có thể chạy tập lệnh dọn dẹp được cung cấp sẵn:

```sh
./remove_rpc.sh
```

Tập lệnh này sẽ tự động loại bỏ toàn bộ mã nguồn, các phụ thuộc, cấu hình docker và các bài kiểm tra liên quan đến RPC, mang lại cho bạn một mẫu Gin REST API hoàn toàn tinh gọn.

## Bắt đầu nhanh

### Phát triển cục bộ (Local development)

```sh
# Khởi chạy Postgres, RabbitMQ, NATS
make compose-up
# Chạy ứng dụng cùng với các thao tác di chuyển cơ sở dữ liệu (migrations)
make run
# Hoặc air
air
# Hoặc debugger (nhấn F5 để dùng với .vscode/launch.json)
```

### Kiểm thử tích hợp (có thể chạy trên hệ thống CI)

```sh
# CSDL, ứng dụng + migrations, kiểm thử tích hợp
make compose-up-integration-test
```

### Triển khai bộ chứa (Docker stack) đẩy đủ với proxy ngược

```sh
make compose-up-all
```

Kiểm tra trạng thái các dịch vụ:

- AMQP RPC:
  - URL: `amqp://guest:guest@127.0.0.1:5672/`
  - Client Exchange: `rpc_client`
  - Server Exchange: `rpc_server`
- NATS RPC:
  - URL: `nats://guest:guest@127.0.0.1:4222/`
  - Server Exchange: `rpc_server`
- REST API:
  - http://app.lvh.me/healthz | http://127.0.0.1:8080/healthz
  - http://app.lvh.me/metrics | http://127.0.0.1:8080/metrics
  - http://app.lvh.me/swagger | http://127.0.0.1:8080/swagger
- gRPC:
  - URL: `tcp://grpc.lvh.me:8081` | `tcp://127.0.0.1:8081`
  - [v1/translation.history.proto](docs/proto/v1/translation.history.proto)
- PostgreSQL:
  - `postgres://user:myAwEsOm3pa55@w0rd@127.0.0.1:5432/db`
- RabbitMQ:
  - http://rabbitmq.lvh.me | http://127.0.0.1:15672
  - Thông tin đăng nhập: `guest` / `guest`
- Trình giám sát NATS (NATS monitoring):
  - http://nats.lvh.me | http://127.0.0.1:8222/
  - Thông tin đăng nhập: `guest` / `guest`

## Cấu trúc dự án

### `cmd/app/main.go`

Tệp cấu hình và khởi tạo hệ thống ghi nhật ký (logger). Sau đó, hàm chính tiếp tục quá trình khởi tạo tại
`internal/app/app.go`.

### `config`

Ứng dụng tuân theo chuẩn 12-factor lưu trữ cấu hình trong các biến môi trường (thường được gọi là `env vars` hoặc `env`). Các biến môi trường
rất dễ dàng thay đổi giữa các đợt phát hành mà không cần sửa đổi mã nguồn. Không giống như các tệp cấu hình thông thường,
rất hiếm khi chúng vô tình bị tải lên kho lưu trữ mã nguồn; và không giống như các cơ chế cấu hình khác (ví dụ: Java System Properties),
chúng là một tiêu chuẩn chung không phụ thuộc vào ngôn ngữ lập trình hay hệ điều hành.

Cấu hình tham khảo: [config.go](config/config.go)

Ví dụ: [.env.example](.env.example)

[docker-compose.yml](docker-compose.yml) sử dụng các biến `env` để cấu hình cho các dịch vụ.

### `docs`

Tài liệu Swagger. Tài liệu này được tự động tạo bởi thư viện [swag](https://github.com/swaggo/swag).
Bạn không cần phải tự sửa đổi tài liệu này theo cách thủ công.

#### `docs/proto`

Các tệp cấu trúc Protobuf. Chúng được sử dụng để tự động sinh mã nguồn (Go) phục vụ cho dịch vụ gRPC.
Các tệp proto này đồng thời cũng được sử dụng để sinh tài liệu cho dịch vụ gRPC.
Bạn không cần phải sửa chửa các thông tin này theo cách thủ công.

### `integration-test`

Dành cho mã kiểm thử tích hợp (Integration tests).
Chúng sẽ được khởi chạy dưới dạng một bộ chứa độc lập (container) chạy song song với bộ chứa của ứng dụng chính.

### `internal/app`

Luôn tồn tại một hàm _Run_ bên trong tệp `app.go`, là chức năng tiếp nối của hàm _main_.

Đây là vị trí khởi tạo tất cả các đối tượng cốt lõi.
Kiến trúc tiêm phụ thuộc (Dependency injection) được thực thi thông qua các hàm khởi tạo "New ..." (vui lòng xem mục Tiêm phụ thuộc).
Kỹ thuật thiết kế này cho phép phân tầng cấu trúc ứng dụng với nguyên tắc [Tiêm Phụ Thuộc (Dependency Injection)](#tiêm-phụ-thuộc-dependency-injection).
Cách tiếp cận này đảm bảo rằng phần logic nghiệp vụ (business logic) hoàn toàn độc lập với các tầng khác trong hệ thống.

Tiếp theo, ứng dụng tiến hành máy chủ ảo và dùng câu lệnh _select_ để đợi tín hiệu tắt nhằm dừng hệ thống một cách an toàn (graceful shutdown).
Trong trường hợp nội dung tệp `app.go` phình to quá mức, bạn có thể thực hiện chia nó ra làm nhiều tệp tin nhỏ.

Đối với trường hợp liên kết một số lượng lớn thành phần ứng dụng, công cụ [wire](https://github.com/google/wire) có thể được ứng dụng.

Tệp `migrate.go` được sử dụng để tiến hành cập nhật hệ thống cơ sở dữ liệu (auto migrations).
Nó chỉ được đưa quy trình biên dịch nếu như tham số cờ truyền vào tương ứng với thẻ ghi chú _migrate_.
Ví dụ:

```sh
go run -tags migrate ./cmd/app
```

### `internal/controller`

Tầng phân tầng lớp kiểm soát các yêu cầu gửi đến (tương ứng với các Controller thuộc mô hình MVC). Mẫu thiết kế minh hoạ 3 loại máy chủ:

- AMQP RPC (tận dụng công nghệ RabbitMQ làm phương tiện truyền tải)
- gRPC (kiến trúc [gRPC](https://grpc.io/) phát triển dựa trên protobuf)
- REST API (Sử dụng kiến trúc không trạng thái với [Gin](https://github.com/gin-gonic/gin) framework)

Tuyến đường định hướng luồng (Routers) ở của máy chủ sở hữu tính kế đồng quy tương đương nhau:

- Mỗi một phần tử điều phối (Handlers) sẽ tập hợp theo hệ sinh thái và phạm vi hoạt động tương thích với ứng dụng.
- Dành cho các phân vùng như vậy, một phần tử định tuyến riêng rẽ sẽ ra mặt cung cấp việc ánh xạ những đường dẫn.
- Đối tượng cấu thành logic nghiệp vụ được chèn vào bên trong cấu trúc định tuyến và sẽ được triệu gọi thông qua bộ xử lý hệ thống.

#### `internal/controller/amqp_rpc`

Quản lý đánh dấu phiên bản dịch vụ RPC.
Theo hình dung ở một phiên bản thế hệ kế cận v2, bạn cần bổ sung riêng cho danh mục này trên nền thư mục mới `amqp_rpc/v2`.
Và bên trong tệp `internal/controller/amqp_rpc/router.go` thực hiện đính kèm dòng mã tương tự sau:

```go
routes := make(map[string]server.CallHandler)

{
    v1.NewTranslationRoutes(routes, t, l)
}

{
    v2.NewTranslationRoutes(routes, t, l)
}
```

#### `internal/controller/grpc`

Quản lý đánh dấu phiên bản dịch vụ gRPC đơn giản.
Tương tự cho thế hệ dịch vụ kết tiếp v2, ta tiến hành cho thư mục mới `grpc/v2` theo như mã nguồn sẵn có.
Chúng ta không thể không nhắc tới là cần đưa kèm theo cấu trúc dữ liệu mô tả mới `v2` với nguyên gốc `docs/proto` bên trên.
Và với tệp tin `internal/controller/grpc/router.go` bổ sung kèm dòng lệnh tương thích:

```go
{
    v1.NewTranslationRoutes(app, t, l)
}

{
    v2.NewTranslationRoutes(app, t, l)
}

reflection.Register(app)
```

#### `internal/controller/nats_rpc`

Đánh dấu phiên bản dịch vụ theo giao thức RPC chuẩn hoá đơn giản.
Ví dụ đối với việc phân bổ phiên bản kế tiếp v2, ta tiến hành đi kèm thư mục `nats_rpc/v2` đi liền với cấu trúc thư viện tương đương.
Và bên trong mục tin cấu thành `internal/controller/nats_rpc/router.go` đính kèm thêm các dòng mã kế tiếp:

```go
routes := make(map[string]server.CallHandler)

{
    v1.NewTranslationRoutes(routes, t, l)
}

{
    v2.NewTranslationRoutes(routes, t, l)
}
```

#### `internal/controller/restapi`

REST versioning chuẩn hoá dễ dàng tiếp cận.
Lấy ví dụ đối chiếu ở phiên bản thay thế thế hệ v2, thực hiện bổ sung một vùng tổ chức mã nguồn `restapi/v2` như thông lệ sẵn có hiện hành.
Và tại mục điều hướng lệnh định tuyến ứng dụng ở `internal/controller/restapi/router.go` khai cấu theo mã lệnh tương ứng:

```go
apiV1Group := app.Group("/v1")
{
	v1.NewTranslationRoutes(apiV1Group, t, l)
}
apiV2Group := app.Group("/v2")
{
	v2.NewTranslationRoutes(apiV2Group, t, l)
}
```

Thay vì tiếp tục triển khai kiến trúc phụ thuộc theo [Gin](https://github.com/gin-gonic/gin), hoàn toàn có được một ứng dụng Http framework mở rộng do bạn tuỳ biến.

Tại thư mục điều hướng tuyến đường `router.go` cùng như phía trên danh sách chức năng bộ phận điều khiển, nhà sản xuất trang bị sẵn các bình luận đặc biệt với mục đích đính kèm với máy chủ Swagger [swag](https://github.com/swaggo/swag) tự động.

### `internal/entity`

Các thực thể liên kết chính của nghiệp vụ (viết tắt là model) có tính cơ động khi sử dụng tuỳ chỉnh ở mọi khu vực logic kiến trúc.
Bên trong đây hệ thống cũng cấp phép với việc sử dụng hàm thực thi đối chiếu, ví dụ, xác thực luồng thông tin kiểm duyệt đầu ra.

### `internal/usecase`

Logic nghiệp vụ (Business logic).

- Các phương thức hoạt động chuyên biệt nhóm dựa vào sự phát triển với hệ sinh thái sử dụng chung của tiện ích
- Mỗi 1 cụm nhóm chuyên biệt sở hữu 1 lớp quản lý cấu trúc độc lập
- Nguyên tắc vận hành cấu tạo trên tương quan Một đối Một: 1 thư mục chứa mã cấu trúc hoạt động với đúng 1 cấu trúc định dạng

Tính mở rộng liên kết tới kho ứng dụng, dịch vụ Web API, dịch vụ RPC hay đa dạng thành phần khác cho lớp logic ứng dụng đều được phép triển khai tiêm mã.
(vui lòng truy cập thông tin thêm tại nhóm [Tiêm phụ thuộc (Dependency Injection)](#tiêm-phụ-thuộc-dependency-injection)).

#### `internal/repo/persistent`

Kho lưu trữ (repository) đại diện cho khái niệm trừu tượng thuộc lưu trữ cơ sở vật chất (có thể là Data) - đóng vai trò là kho giao tiếp chính xác đến lớp hệ thống hoạt động logic chuyên trách.

#### `internal/repo/cache`

Đóng vai trò như một môi trường trừu tượng cung cấp tính linh hoạt trên bộ nhớ đệm ứng dụng dành cho môi trường logic chuyên trách.
Xây dựng lớp liên kết (adapter pattern) hỗ trợ các cấu tạo bộ nhớ đệm đa nền tảng hiện nay (ví dụ: máy chủ Otter tại cấu trúc `pkg/cache`).

#### `internal/repo/webapi`

Trình biên diễn trên cơ chế phi tập trung cho môi trường trang Web API kết hợp lớp logic chung ứng dụng.
Lấy ví dụ cơ bản, bộ quản lý này đóng vai logic nền hướng ra một ứng dụng vi mạng liên kết hệ đa bên có cung cấp đường dẫn REST API nội liên.
Tên hệ thống gói mở đầu do các thiết định cơ sở để đưa ra những danh mục tối ưu tính độc lấp theo tuỳ mục đích.

### `pkg/rabbitmq`

Chuẩn liên kết cấu trúc RabbitMQ RPC:

- Không sử dụng cấu trúc định tuyến phức hợp nằm sâu ở RabbitMQ
- Khuyến nghị sử dụng luồng truyền tải Exchange kiểu `fanout`, kết hợp 1 chuỗi độc lập exclusive queue nhằm đạt được lượng tài nguyên công suất tối đa
- Hệ thống cơ chế cung cấp khả năng Kết Nối Lại tự thiết lập cho ứng dụng để tránh trạng thái Mất kết nối liên tục

### `pkg/cache`

Lớp cung cấp dữ liệu theo chuẩn phân giải cao về độ nhớ đệm hệ thống [Otter](https://github.com/maypok86/otter):

- Quản lý kho bộ nhớ đệm chuẩn tốc độ cao (high-performance) ở cơ sở nội tại nền và quản lý cấp nhật hoạt động bảo mật tài nguyên theo Thread-Safe nội tệp định hình.
- Đi kèm tính khả kết xuất hạn mức hệ thống xoay vòng thông tin như Cạn kiệt Thời hạn TTL và quản lý hạn chế tốn hao không gian nhớ lưu tối đa
- Có đi kèm định hình tích luỹ đơn nhiệm gọi hàm `GetOrLoad` trong mô hình Singleflight đính kèm hỗ trợ quy chuẩn chống xả trễ thời điểm và làm chậm bộ xử lý do tràn độ trễ nhớ trong khoảng thời gian xác quyết.

## Tiêm phụ thuộc (Dependency Injection)

Theo như tiêu chí chia tách độ ràng buộc phụ thuộc trên những tệp tin nội hệ thống ra các môi trường nền thuộc thư mục bên ngoài độc lập, tính năng truyền biến kiến thiết Tiêm phụ thuộc (DI) luôn nhận được giá trị cao nhất.

Theo đơn cử, đi cùng cấu trúc hàm dựng (constructor) như chữ khởi đầu `New`, toàn bộ chu trình có xu hướng đẩy các tuỳ kết cấu của luồng thư viện đi qua tầng này tới môi trường logic thuần thực chiến.
Khối tài liệu thiết lập đem lại giá trị bảo đảm hoạt động độc quyền và hệ kết nối cho Logic thông tin (trên độ mở tuỳ động và hỗ trợ vận dụng chuyển đổi tài liệu dễ cấu tạo).
Các tính năng hiện nay còn cho ra phương hướng bổ sung những tuỳ biến mở khoá kết nối giao thức, tuy nhiên toàn bộ tiến trình không tác động hoặc buộc can thiệp đến hệ thực thể gói `usecase`.

```go
package usecase

import (
// Không bao gồm khai báo phụ thuộc!
)

type Repository interface {
	Get()
}

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) Do() {
	uc.repo.Get()
}
```

Cách thiết lập cơ bản này tạo cơ chế dễ tiếp vận các tính định chuẩn từ môi trường mô phỏng lập trình Mock (tính đến như [mockery](https://github.com/vektra/mockery)) và kết xuất cho một kho tiện nghi thực thi các chu trình kiểm đoạn thông tin đơn giản hơn (unit tests).

> Chúng ta không ràng buộc cứng vào các bản triển khai trực tiếp, cho phép khả năng thay thế linh hoạt cấu phần này với cấu phần khác mà không gây ảnh hưởng lẫn nhau.
> Nếu cấu phần ứng dụng mới áp kết quy chuẩn giao thức chung thông qua (interface), thì bạn không cần thực hiện tinh chỉnh tái quy hoạch nào trong cấu trúc tệp mã business logic.

## Kiến trúc Sạch (Clean Architecture)

### Ý tưởng trọng tâm

Các lập trình viên thường nhận diện cấu trúc ứng dụng tối ưu nhất sau bước hoàn tất và thiết kế tương đối phần lớn danh giới những mã luồng quan trọng trong hệ thống.

> Một cấu trúc thiết kế ứng dụng đạt mức tối đa hóa hiệu quả phải cho ra thời gian tạo quyết nghị bị trì hoãn càng lùi lại về thời khắc đoạn cuối khi bạn thực thi thì nhận diện cơ hội để hoàn thiện chuẩn định của thông tin còn bao giờ hết.

### Nguyên lý trọng tâm then chốt

Sức nặng của thuật toán đảo nghịch sự phụ thuộc mang danh Nguyên lý Đảo ngược Phụ thuộc ("Dependency Inversion" một yếu tố chữ D lấy từ quy trình S.O.L.I.D đại chúng) cũng dựa vào những tính năng cốt lỗi bên trên từ "Tiêm phụ thuộc".
Sự vận động hướng liên thông từ kiến trúc các vùng ở tầm vĩ mô và đi sâu sắc hơn vào vòng cấu tạo nền vi mô ở lõi sâu thẳm.
Sở dĩ ứng dụng phát lệnh tạo liên hệ thành công là ở mặt thực hành, cấu trúc nghiệp vụ thông minh cũng được ưu thích để bảo quản hoạt động tách biệt tự do ra các vùng biên cấu tạo ra môi trường ngoại vi khép mở trong tổng khu.

Do vậy thiết kế, mô hình phần mềm được định biên ra hai ranh giới phân tầng: tầng ngoại vi bên ngoài và giới tuyến nội khu trong cốt lõi:

1. **Vùng Logic nghiệp vụ thông tin nội tại** (bộ cung ứng thư viện theo chuẩn của nền tảng ngôn ngữ Go).
2. **Khu vực mở rộng hạ tầng** (ví dụ thông quản Cơ sở dữ liệu, môi giới cổng phân vùng thông tin thư điện gửi trao tin nhắn, đến từng danh mục cấu phần bên khu vực bao gồm hoặc kiến trình xây dựng ngoài phạm vi của kiến trúc này).

![Kiến trúc hiện đại theo Clean Architecture](docs/img/layers-1.png)

**Lớp lõi mang yếu tố trong cùng** gắn chặt chẽ vào quy trình quy chuẩn kiến trúc sạch và khép độc lập. Nó buộc phải:

- Không sử dụng các gọi gói mã cài đè thông qua cấu trúc thư mục thuộc vùng phân chia ngoại vi lớp ngoài.
- Dành sự chuyên biệt dùng cấu phần thư viện tiêu chuẩn hỗ trợ trên cùng công nghệ nguyên bản.
- Áp quy luật gọi các cơ cấu cấu trúc ngoài vùng và thực thi trên cấu phần Interface chung (!).

Vùng thực thể phân chuyên nghiệp vụ là độc lập không nắm các kiến thức bên hệ máy quản dữ liệu Postgres và hệ phân mảng về ứng dụng thiết lập công khai.
Cấu phân bộ giao thức sẽ mang riêng 1 kết nối bộ lập trình đa tương quan phục vụ vận hành cùng dữ liệu máy CSDL với API chuẩn qua cách tính độ ẩn tàng (trừu tượng).

**Tuyến phòng thủ vùng phía ngoài biên** bị hạn chế do nhiều yêu cầu ràng buộc:

- Mọi module thành tố phân chia giới từ cấu tạo của không gian vùng đều hoạt động và giàn quản lí cách ngăn cho tương quan giao nhau. Một thiết định để mở tính kết nối ở 2 loại thành tố đó như thế nào? Không dựa vào hàm mở gọi nối biên ngay được, chỉ vận dụng tính trừu cấp thông qua nội vòng trong cơ sở mã thuộc Lõi Logic hoạt động của chính bản thông tin cấu thành ứng dụng logic nghiệp vụ.
- Tất thảy yêu cầu hướng tâm nội vùng là do phương diện thông vận kết qua hệ Giao Thức Lập Trình Interface quy chuẩn (!).
- Việc chia dữ liệu để di chuyển có cấu tạo định chuẩn thông nhất tạo tính linh cơ hỗ trợ thiết lập Logic theo cách nhanh chập (`internal/entity`).

Lấy định dạng, lúc mà công đoạn mở kết cấu lưu hệ thông kho DB và khởi đầu gọi từ HTTP HTTP Endpoint (với đối tác controller tương thức):
HTTP và DB hai cấu phần được bố trí ở ngoại khu nên có sự thiết thiết phải thiết tập không hề biết ranh giữa hoạt động của nhau.
Do những ràng hệ cách nhau đó được gọi gián kết nối qua một quy cơ mang hình bóng qua các vùng `usecase` (là kiến trúc logic định phận):

```
    HTTP > usecase
           usecase > repository (Lõi Postgres chứa DB)
           usecase < repository (Lõi Postgres chứa DB)
    HTTP < usecase
```

Bảo vận cho dấu tự biên của > và biểu trưng nhỏ hơn < có điểm hình thức giao quy cắt ngang từ đường lằn tầng tuyến phân biên.
Qua sự phản ánh được dẫn dắt ở biểu đồ phân mục tổng khu:

![Biểu định hướng hoạt định](docs/img/example-http-db.png)

Để phức biên các logic đa kết liên qua hệ thống cấu trình độ khó lớn:

```
    HTTP > usecase
           usecase > repository
           usecase < repository
           usecase > webapi
           usecase < webapi
           usecase > RPC
           usecase < RPC
           usecase > repository
           usecase < repository
    HTTP < usecase
```

### Tiêu cấu của các Lớp Layering

![Ví dụ](docs/img/layers-2.png)

### Những cụm ngôn từ thuộc mô hình chuẩn Clean Architecture

- **Thành Cố Các Thực Thể** (mã `internal/entity`) nằm ở cơ chế nền thiết lập trung tâm vận quan với dữ liệu có tính cách kết tập vào tất thảy những thiết định có sẵn nằm xung quanh vùng liên tục hệ phỏng lập nghiệp. Cụ thể có thể ứng dụng trong quá trình quy chuyển xuyên thấu hệ môi trường toàn tuyến. Chúng độc lập vắng bỏ hoàn tính kết liên vĩnh viễn khỏi các kho tài dữ liệu ngoài vùng máy lưu hệ cơ cấu đường giao kết mạng.
- **Dữ liệu Hình Khối Dữ Liệu Hoá Cơ Sở DB Models** (với đường liên kết tạo ở `internal/repo/persistent/sqlc`) là đoạn nền mã do tự lập hình kiến thiết và minh chứng những bản thể mô phẳng kết nối với không gian sơ dồ kho dữ liệu phân vùng. Bộ thành phận nằm tại các hệ kín (thuộc bảo vận quyền khép vòng không gian repo layer - repository kho cơ sở).
- **Trình Kết Chuyển Ánh Xạ Mapping:** Nằm khu vực cơ phận vùng bộ liên hợp kho lưu repo được chuyên giao trong luồng tạo ra trách nhiệm phân tính của việc thực hiện thao phân dịch cơ thể quy chế biến ("mapping") khối hệ cơ chế mô phỏng hệ csdl sang định dạng của thực khuyết bên phân nhánh kinh danh thông vận ("entities"). Mục đích tối ưu là nhằm khẳng duy tính hoạt động toàn chu không để cho một tuỳ thể tinh chỉnh theo các hướng phân loại của biểu đồ liên dữ liệu schema tàn kích làm phân phả và chỉnh lệch biến dạng thành quy vào các quy luật chung của nghiệp hành chính tại bản quy định phân cấp nghiệp vụ ứng logic.

- **Tiến trình Hành Vy Kinh Phân Use Cases** đưa hệ luật lưu hành cho công cuộc thiết chế quy chế định quy tại bên trong tệp `internal/usecase`.

Tại lớp quy trình logic thông vận đang chuyên hoạt hoá tương duy theo cách tiếp tác với nhau mang những mỹ ngữ có gọi bằng cái tên lớp phân tầng _cơ sở kiến tầng cấu kết hạ tầng_ (_infrastructure_ layer).
Được triển tác như kho vùng dữ kho `internal/usecase/repo`, bộ thiết mở từ vùng nội ngoại ranh như giao liên webapi `internal/usecase/webapi`, tới tất thảy các gói cơ phận hệ thống (các dạng pkg chung mở ra tính điều khiển các trạm ngoại kho khác) tính tới với chu dịch phân nhỏ liên hệ microservices khác nhau.
Cho qua một mẫu đồ nguyên tính template có đính ở tệp hướng, khối các thành tố nằm theo cơ bản cơ sở giao cấu liên tính (_infrastructure_) được uốn lưu bên trong khung ngầm ranh bộ khung kho ứng mục `internal/usecase`.

Mục đích thiết biên quy đặt lên quyền lợi mở lối cho cửa thông quy của thiết định cửa ứng thông quan hướng ra thì định theo ngôn ngữ tự nguyên ở tuỳ theo thẩm quyền quyết đinh:

- hệ kiểm soát (áp ứng cho thiết thiết bản quy này theo controller)
- khu gửi thông ứng delivery
- khu giao trạm uỷ đường truyền phân lưu giao vận transport
- cửa giao điều khu phân tính (gateways)
- đầu mục điểm nối khởi phát (entrypoints)
- bộ liên giao thông tiểu gốc trọng nguyên primary
- cửa điều dẫn vào phân cấp input

### Cấu tầng đệm Phụ trợ cấu tầng (Additional layers)

Đối chuẩn cho dòng Cấu trình Kiến tạo cổ điển thông chuẩn phân theo Kiến Trúc Đồ theo Nguyên Lý "Kiến Trúc Tĩnh Khung Bạch Rõ" - [kiến cơ từ chuyên bài The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) được lập phát từ quá trình khởi kết của một Cấu Khối ứng dạng Khối Cấu Đại Kết Đặc Tĩnh Lớn Lập Liên Trụ Toàn Khung monolithic qua 4 Lớp cấu quy phân đệ tầng.

Dịch ra phương bản hướng cổ điển chuẩn nhất cho sự phân giải tầng ngoại được đem tác cắt nửa thành cấu tạo 2 thành phận, có thêm việc sở hữu và thiết hành khả chuyển hoá hướng đảo lại tuân theo định hình đảo ngược thiết chuẩn ngoại hệ thống cũng có mang quy phân nội dẫn. Tạo giao thức định tuyến kết tương qua lớp bề mặt phương giao truyền lệnh mặt phẳng (dùng cơ sở Interfaces).

Tại vị trí điểm Lõi tầng trong Nội Khuyết Lớp Trong (khi logic có cấu phần phức giải và cộm rườm) cũng đã cắt tác ra cấu biên ra định chia 2 luồng nửa ranh, qua đính áp thêm phân ly giao thức.

---

Kết lại về cách hướng điều tính vận quy với bộ chức có sự đố khối với khả kết nối vĩ hệ lớn cần ra sức phân quyền điều bộ thiết lớp bổ uốn giao luồng cơ bộ định bổ sung tuỳ phát tầng chuyên hệ để mở thiết luồng quản ranh liên động.
Tuy nhiên, cấu trúc này là hệ thiết bị tính đặc lưu do đó các nhà khoa bộ triển định nên cần "CHỈ BỔ SUNG KHÔNG GIAN CẤU TẦNG" tính theo trường khả biến mà khi và chỉ khi lúc yêu sách từ hệ uốn định mang quy cấp thiết THẬT CỰC KỲ KHẨN TRONG ƯỚC PHÁT (Really Necessary Cần Nổi Lúc Quyết Khởi) Lớp kiến thành mới được thiết thiết lập quy chế định chuẩn (như do mã tại thông kỳ lúc đó đã khống bám không uốn theo).

### Các quy cách giao cấu thiết định mở rộng đa cấu (Alternative approaches)

Trừ qua cái tên thuộc quy theo định "Kiến trúc hệ Tĩnh hệ Thuần Cấu Khởi" (Clean Architecture), mô cấu thuộc thiết "Khối Hình phân luỹ theo Kiểu Onion - Khối Hành Trưởng Khối" (Onion architecture) cho với khối hướng qua hình "Lục Cấu Định Giải Biểu Lục Giác - Hexagonal architecture" (Hay là định Cổng phân uốn Mở và Lô Mô Đun Đầu Giải Biểu: Ports and adapters) do được đính tính cũng song chiếu theo như chuẩn chung nhất.
Hiệp 2 tính thể mô cấu phương định dạng có thuộc và xuất cấu chuẩn theo điểm Lấy từ quy phân Tĩnh Thuần Dấu Ngược (Dependency Inversion).
Để nhìn nhận sâu trong kết hình khung hình "Khối bộ khung chia Cửa truyền cùng mô phân cổng Mở giải giao" _Ports and adapters_ có chung các tương thông cực điểm theo _Clean architecture_, và khác đi nằm duy tại vấn của điểm thuật biểu ngôn xưng pháp gọi (terminologies).

## Một số cấu dự có thành phận cấu đồng hệ tương đẳng

- [https://github.com/evrone/go-clean-template](https://github.com/evrone/go-clean-template)
- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## Trích tham hướng nguồn thông kết thông tham khảo hữu bổ (Useful links)

- [Bài báo The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Nguyên thức của Twelve factors](https://12factor.net/)
