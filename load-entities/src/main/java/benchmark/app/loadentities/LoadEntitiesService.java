package benchmark.app.loadentities;

import com.googlecode.objectify.cmd.QueryKeys;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.Stream;
import java.util.stream.StreamSupport;

import static com.googlecode.objectify.ObjectifyService.ofy;

import com.googlecode.objectify.ObjectifyService;

@Service
public class LoadEntitiesService {

    public LoadEntitiesService() {
        ObjectifyService.register(DummyEntity.class);
    }

    public void generateEntities(int numberOfEntities) {
        Stream<Integer> intStr = IntStream.range(0, numberOfEntities).boxed();

        List<DummyEntity> entities = intStr.map(val -> DummyEntity.createNew()).collect(Collectors.toList());
        ofy().save().entities(entities).now();
    }

    public List<String> getLoadedKeys() {
        QueryKeys<DummyEntity> keys = ofy().load().type(DummyEntity.class).keys();

        return StreamSupport.stream(keys.spliterator(),false).map(key -> key.toWebSafeString()).collect(Collectors.toList());
    }
}