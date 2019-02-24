package benchmark.spring.maincontroller;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

@RestController()
public class MainController {
    private static Log logger = LogFactory.getLog(MainController.class);
    private MainService mainService;

    @Autowired
    public MainController(MainService mainService) {
        this.mainService = mainService;
    }

    @GetMapping("entity/{id}")
    @ResponseBody
    public DummyEntity fetchEntity(@PathVariable("id") String id) {
        logger.info("Request to fetch entity " + id);
        return mainService.fetchById(id);
    }
}
