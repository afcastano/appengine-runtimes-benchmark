package benchmark.app.loadentities;

import com.googlecode.objectify.annotation.Entity;
import com.googlecode.objectify.annotation.Id;

import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;

@Entity
public class DummyEntity {
    @Id
    private String id;
    private String random1;
    private Integer random2;

    public DummyEntity() {}

    public static DummyEntity createNew() {
        DummyEntity dummyEntity = new DummyEntity();

        dummyEntity.id = RandomStringUtils.randomAlphanumeric(8);
        dummyEntity.random1 = RandomStringUtils.randomAlphanumeric(8);
        dummyEntity.random2 = RandomUtils.nextInt();

        return dummyEntity;
    }
}