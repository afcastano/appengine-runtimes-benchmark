package benchmark.thundr.fetchentities;
import com.googlecode.objectify.annotation.Entity;
import com.googlecode.objectify.annotation.Id;
import com.googlecode.objectify.annotation.Index;

@Entity
public class DummyEntity {
    @Id
    private String id;
    private String random1;
    @Index
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
