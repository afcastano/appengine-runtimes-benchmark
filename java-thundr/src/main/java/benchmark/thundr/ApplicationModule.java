package benchmark.thundr;

import benchmark.thundr.fetchentities.DummyEntity;
import benchmark.thundr.fetchentities.FetchEntitiesController;
import benchmark.thundr.fetchentities.FetchEntitiesService;
import com.googlecode.objectify.ObjectifyService;
import com.threewks.thundr.gae.GaeModule;
import com.threewks.thundr.gae.objectify.ObjectifyModule;
import com.threewks.thundr.injection.BaseModule;
import com.threewks.thundr.injection.UpdatableInjectionContext;
import com.threewks.thundr.module.DependencyRegistry;
import com.threewks.thundr.route.Router;

public class ApplicationModule extends BaseModule {
	public static class Route {
		public static final String Health = "health";
		public static final String FetchEntity = "fetchEntity";
	}

	@Override
	public void requires(DependencyRegistry dependencyRegistry) {
		super.requires(dependencyRegistry);
		dependencyRegistry.addDependency(GaeModule.class);
		dependencyRegistry.addDependency(ObjectifyModule.class);
	}

	@Override
    public void initialise(UpdatableInjectionContext injectionContext) {
        super.initialise(injectionContext);
        configureObjectify();
    }

	@Override
	public void configure(UpdatableInjectionContext injectionContext) {
		super.configure(injectionContext);
		injectionContext.inject(FetchEntitiesService.class).as(FetchEntitiesService.class);
	}

	@Override
	public void start(UpdatableInjectionContext injectionContext) {
		super.start(injectionContext);
		Router router = injectionContext.get(Router.class);
		addRoutes(router);
	}

	private void addRoutes(Router router) {
		router.get("/", FetchEntitiesController.class, "sayHello", Route.Health);
		router.get("/entity/{id}", FetchEntitiesController.class, "fetchEntity");
		router.get("/entity/greaterThan/{index}", FetchEntitiesController.class, "fetchEntities");
	}

	private void configureObjectify() {
		ObjectifyService.register(DummyEntity.class);
	}
}
