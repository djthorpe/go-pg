// The `test` package supports running unit tests with containerized supporting
// services. This is useful for testing that require external services, such as
// databases, message brokers, etc. The `test` package provides a simple API to
// create and manage containers.
//
// For example, in order to test postgres integration tests, use the following
// boilerplate for your tests:
//
//		// Global variable which will hold the connection
//		var conn test.Conn
//
//	 // Start up a container and return the connection
//		func TestMain(m *testing.M) {
//						test.Main(m, &conn)
//					}
//
//			     // Run a test which pings the database
//					func Test_Pool_001(t *testing.T) {
//						assert := assert.New(t)
//						conn := conn.Begin(t)
//						defer conn.Close()
//
//						// Ping the database
//						assert.NoError(conn.Ping(context.Background()))
//					}
package test
