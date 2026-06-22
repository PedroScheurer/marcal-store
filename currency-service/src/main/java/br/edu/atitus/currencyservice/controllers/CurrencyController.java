package br.edu.atitus.currencyservice.controllers;

import br.edu.atitus.currencyservice.dtos.CurrencyDTO;
import br.edu.atitus.currencyservice.entities.CurrencyEntity;
import br.edu.atitus.currencyservice.infrastructure.exceptions.CurrencyNotFoundException;
import br.edu.atitus.currencyservice.repositories.CurrencyRepository;
import br.edu.atitus.currencyservice.services.CurrencyService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/currency")
public class CurrencyController {

    private final CurrencyService currencyService;

    public CurrencyController(CurrencyService currencyService) {
        this.currencyService = currencyService;
    }

    @GetMapping(path = "/convert",params = {"source", "target"})
    public ResponseEntity<CurrencyDTO> findBySourceAndTarget(@RequestParam String source,
                                                            @RequestParam String target){

        CurrencyDTO dto = currencyService.findBySourceAndTarget(source,target);

        return ResponseEntity.ok().body(dto);
    }
}
