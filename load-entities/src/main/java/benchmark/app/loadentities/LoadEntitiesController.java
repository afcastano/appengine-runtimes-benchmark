package benchmark.app.loadentities;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController("/load")
public class LoadEntitiesController {

    private LoadEntitiesService service;

    @Autowired
    public LoadEntitiesController(LoadEntitiesService service) {
        this.service = service;
    };

    @PostMapping
    public String createEntities() {
        service.generateEntities(200);
        return "200 entities created";
    }

    @GetMapping
    public List<String> getKeys() {
        return service.getLoadedKeys();
    }
    
}