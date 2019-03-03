package benchmark.thundr.fetchentities;

import com.threewks.thundr.view.json.JsonView;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import java.util.List;

public class FetchEntitiesController {

    private static Log logger = LogFactory.getLog(FetchEntitiesController.class);

    private FetchEntitiesService service;

    public FetchEntitiesController(FetchEntitiesService service) {
        this.service = service;
    }

    public JsonView sayHello() {
        return new JsonView("Thundr service up and running!!!");
    }

    public JsonView fetchEntity(String id) {
        logger.info("Request to fetch entity " + id);
        return new JsonView(service.fetchById(id));
    }

    public JsonView fetchEntities(Integer index) {
        logger.info("Request to query entity greater than " + index);
        List<DummyEntity> foundEntities = service.queryGreaterThanIndex(index);
        logger.info("Found " + foundEntities.size() + " entities");
        return new JsonView(foundEntities);
    }
}
