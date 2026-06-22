package dev.pedroscheurer.greetingservice.controllers;

import dev.pedroscheurer.greetingservice.configs.GreetingConfig;
import dev.pedroscheurer.greetingservice.dtos.GreetingDTO;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/greeting")
public class GreetingController {

    private final GreetingConfig config;

    public GreetingController(GreetingConfig config) {
        this.config = config;
    }

    @GetMapping
    public String getGreetingParam(
            @RequestParam(required = false) String name) {
        if (isNameInvalid(name)) {
            name = config.getDefaultName();
        }
        return formatedReturn(name);
    }

    @GetMapping({ "/{name}"})
    public String getGreetingVariable(
            @PathVariable String name) {
        if (isNameInvalid(name)) {
            name = config.getDefaultName();
        }
        return formatedReturn(name);
    }

    @PostMapping
    public String getGreeting(@RequestBody GreetingDTO dto){
        String name = dto.name();
        if (isNameInvalid(name)) {
            name = config.getDefaultName();
        }
        return formatedReturn(name);
    }

    private boolean isNameInvalid(String name){
        return name == null || name.isEmpty();
    }

    private String formatedReturn(String name){
        return String.format("%s %s!!!", config.getGreeting(), name);
    }
}