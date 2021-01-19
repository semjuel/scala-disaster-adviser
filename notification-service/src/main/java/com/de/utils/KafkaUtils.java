package com.de.utils;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.IntegerDeserializer;
import org.apache.kafka.common.serialization.LongDeserializer;
import org.apache.kafka.common.serialization.StringDeserializer;
import reactor.kafka.receiver.ReceiverOptions;

import java.time.Duration;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

public class KafkaUtils {
    private KafkaUtils() {

    }

    public static ReceiverOptions<Integer, String> receiverOptions(String bootstrapServers,
                                                                   Collection<String> topics,
                                                                   String clientId,
                                                                   String groupId,
                                                                   Duration commitInterval) {
        Map<String, Object> props = new HashMap<>();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        props.put(ConsumerConfig.CLIENT_ID_CONFIG, clientId);
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
        props.put(ConsumerConfig.GROUP_ID_CONFIG, groupId);

        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        return ReceiverOptions.<Integer, String>create(props).subscription(topics).commitInterval(commitInterval);
    }
}
