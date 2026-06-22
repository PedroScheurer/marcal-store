package br.edu.atitus.order_service.services;

import br.edu.atitus.order_service.dtos.OrderDTO;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import br.edu.atitus.order_service.clients.CurrencyClient;
import br.edu.atitus.order_service.clients.CurrencyResponse;
import br.edu.atitus.order_service.clients.ProductClient;
import br.edu.atitus.order_service.clients.ProductResponse;
import br.edu.atitus.order_service.entities.OrderEntity;
import br.edu.atitus.order_service.entities.OrderItemEntity;
import br.edu.atitus.order_service.repositories.OrderRepository;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Service
public class OrderService {

    private final OrderRepository orderRepository;
    private final ProductClient productClient;
    private final CurrencyClient currencyClient;

    public OrderService(OrderRepository orderRepository, ProductClient productClient, CurrencyClient currencyClient) {
        this.orderRepository = orderRepository;
        this.productClient = productClient;
		this.currencyClient = currencyClient;
    }

    public OrderEntity createOrder(OrderDTO orderDTO, Long userId) {
        OrderEntity order = new OrderEntity();
        order.setOrderDate(LocalDateTime.now());
        order.setCustomerId(userId);

        double totalPrice = 0.0;
        double totalConvertedPrice = 0.0;
        String targetCurrency = "USD";

        List<OrderItemEntity> items = new ArrayList<>();

        for (var itemDTO : orderDTO.items()) {
            OrderItemEntity item = new OrderItemEntity();
            item.setProductId(itemDTO.productId());
            item.setQuantity(itemDTO.quantity());

            ProductResponse product = productClient.getProductById(itemDTO.productId());
            item.setPriceAtPurchase(product.price());
            item.setCurrencyAtPurchase(product.currency());
            item.setProduct(product);

            totalPrice += product.price() * itemDTO.quantity();

            CurrencyResponse currencyResponse = currencyClient.getCurrency(product.currency(), targetCurrency);
            double convertedPrice = product.price() * currencyResponse.getConversionRate();
            item.setConvertedPriceAtPruchase(convertedPrice);

            totalConvertedPrice += convertedPrice * itemDTO.quantity();

            item.setOrder(order);
            items.add(item);
        }

        order.setItems(items);
        order.setTotalPrice(totalPrice);
        order.setTotalConvertedPrice(totalConvertedPrice);

        return orderRepository.save(order);
    }

    public Page<OrderEntity> findOrdersByCustomerId(Long customerId, String targetCurrency, Pageable pageable) {
    	Page<OrderEntity> orders = orderRepository.findByCustomerId(customerId, pageable);
    
    	
    	for (OrderEntity order : orders) {
    		double totalPrice = 0.0;
        	double totalConvertedPrice = 0.0;
        
            for (OrderItemEntity item : order.getItems()) {
                ProductResponse product = productClient.getProductById(item.getProductId());
                item.setProduct(product);
                totalPrice += item.getPriceAtPurchase() * item.getQuantity();
                
                CurrencyResponse currencyResponse = currencyClient.getCurrency(item.getCurrencyAtPurchase(), targetCurrency);
                item.setConvertedPriceAtPruchase(item.getPriceAtPurchase() * currencyResponse.getConversionRate());
                totalConvertedPrice += item.getConvertedPriceAtPruchase() * item.getQuantity();
            }
            order.setTotalPrice(totalPrice);
            order.setTotalConvertedPrice(totalConvertedPrice);
        }
        return orders;
    }
}
