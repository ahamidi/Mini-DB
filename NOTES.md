### Notes

* DB-wide Locking - I included a DB-wide lock (mutex) in order to synchronize the creation of keys.
* Record Locking - I used per record locks to achieve the "waiting" required by the specs.

### 3rd Party Packages

* [Negroni](https://github.com/codegangsta/negroni) - Middleware. Provides niceties such as recovery.
* [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux)
* [Testify](https://github.com/stretchr/testify) - Convenience wrappers for assertions.
