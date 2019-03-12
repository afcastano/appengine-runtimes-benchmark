package benchmark.spring.maincontroller;

import com.googlecode.objectify.ObjectifyService;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.stereotype.Service;

import java.util.List;

import static com.googlecode.objectify.ObjectifyService.ofy;

@Service
public class MainService {

    private static Log logger = LogFactory.getLog(MainService.class);


    public MainService() {
        ObjectifyService.register(DummyEntity.class);
    }

    public DummyEntity fetchById(String id) {
        logger.info("Fetching entity " + id);
        DummyEntity found = ofy().cache(false).load().type(DummyEntity.class).id(id).now();

        if (found == null) {
            logger.info("Entity not found");
        } else {
            logger.info("Found entity " + found.getId());
        }

        return found;
    }

    public List<DummyEntity> queryGreaterThanIndex(Integer random2) {
        logger.info("Querying entities greater than " + random2);
        List<DummyEntity> entities = ofy().cache(false).load().type(DummyEntity.class)
                .filter("random2 >=", random2).filter("random2 <", random2 + 10000).limit(10).list();

        return entities;
    }


}
