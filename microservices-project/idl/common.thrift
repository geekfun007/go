namespace go common

// 基础响应
struct BaseResponse {
    1: required i32 code
    2: required string message
}

// 分页请求
struct PageRequest {
    1: optional i32 page = 1
    2: optional i32 page_size = 10
}

// 分页响应
struct PageInfo {
    1: required i32 page
    2: required i32 page_size
    3: required i64 total
}
