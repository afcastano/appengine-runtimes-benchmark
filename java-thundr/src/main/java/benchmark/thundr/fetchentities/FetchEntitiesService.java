package benchmark.thundr.fetchentities;

import com.googlecode.objectify.Key;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import static com.googlecode.objectify.ObjectifyService.ofy;

public class FetchEntitiesService {
    private static Log logger = LogFactory.getLog(FetchEntitiesService.class);

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
