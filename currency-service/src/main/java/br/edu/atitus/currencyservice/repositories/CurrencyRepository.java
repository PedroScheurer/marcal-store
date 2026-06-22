package br.edu.atitus.currencyservice.repositories;

import br.edu.atitus.currencyservice.entities.CurrencyEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface CurrencyRepository extends JpaRepository<CurrencyEntity, Long> {

    Optional<CurrencyEntity> findBySourceCurrencyAndTargetCurrency(String sourceCurrency,
                                                                   String targetCurrency);
}
