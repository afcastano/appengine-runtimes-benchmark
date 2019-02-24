package benchmark.spring.maincontroller;

import com.googlecode.objectify.annotation.Entity;
import com.googlecode.objectify.annotation.Id;

@Entity
public class DummyEntity {
    @Id
    private String id;
    private String random1;
    private Integer random2;

    public DummyEntity() {}

    public String getId() {
        return id;
    }

    public String getRandom1() {
        return random1;
    }

    public Integer getRandom2() {
        return random2;
    }
}
