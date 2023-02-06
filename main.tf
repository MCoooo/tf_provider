/* resource "example_server" "my-server-name" {
	uuid_content = "8"
} */

resource "example_server" "pod1" {
    pod {
        name = "apples"
        build {
            ram = 9
            disk = 10199
        }
    }
}





/* data "example_server" "uuid_content" {
    get_len = 1
    uuid_content = "1"
} */

/* output "my-data" {
  value = data.example_server.uuid_content
} */

/* resource "example_server" "file1" {
	uuid_count = "8"
} */
