package benchmark.spring.maincontroller;

import com.googlecode.objectify.Key;
import com.googlecode.objectify.ObjectifyService;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.stereotype.Service;

import static com.googlecode.objectify.ObjectifyService.ofy;

@Service
public class MainService {

    private static Log logger = LogFactory.getLog(MainService.class);


    public MainService() {
        ObjectifyService.register(DummyEntity.class);
    }

    public DummyEntity fetchById(String id) {
        logger.info("Fetching entity " + id);
        DummyEntity found = (DummyEntity) ofy().cache(false).load().key(Key.create(id)).now();

        if (found == null) {
            logger.info("Entity not found");
        } else {
            logger.info("Found entity " + found.getId());
        }

        return found;
    }


}
